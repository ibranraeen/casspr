package mermaid

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// Server represents the local HTTP and WebSocket server for interactive
// Mermaid diagram viewing and editing.
type Server struct {
	listener net.Listener
	mux      *http.ServeMux
	code     string
	mu       sync.Mutex
	onChange func(string)
	onSync   func(string)
	conns    map[*websocket.Conn]bool
	server   *http.Server
}

// WSMessage represents the schema of messages exchanged over the WebSocket connection.
type WSMessage struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

// NewServer creates a new local Mermaid editor server.
func NewServer(initialCode string, onChange func(string), onSync func(string)) *Server {
	return &Server{
		mux:      http.NewServeMux(),
		code:     initialCode,
		onChange: onChange,
		onSync:   onSync,
		conns:    make(map[*websocket.Conn]bool),
	}
}

// Start listens on a random loopback address and starts the HTTP server.
// It returns the listener's local address (host:port) that was assigned.
func (s *Server) Start() (string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("failed to listen on loopback: %w", err)
	}
	s.listener = ln
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.Handle("GET /ws", websocket.Handler(s.handleWS))

	s.server = &http.Server{
		Handler: s.mux,
	}

	go func() {
		if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			slog.Error("Mermaid server failed", "error", err)
		}
	}()

	return ln.Addr().String(), nil
}

// Addr returns the assigned host:port address the server is listening on.
func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}

// Stop shuts down the running HTTP server.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		conn.Close()
	}

	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// UpdateCode updates the diagram code in the server and broadcasts the update
// to all connected editor clients.
func (s *Server) UpdateCode(code string) {
	s.mu.Lock()
	s.code = code
	s.mu.Unlock()

	s.broadcast(WSMessage{
		Type: "update",
		Code: code,
	})
}

// broadcast sends a message to all active WebSocket connections.
func (s *Server) broadcast(msg WSMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal WebSocket broadcast", "error", err)
		return
	}

	for conn := range s.conns {
		_, err := conn.Write(payload)
		if err != nil {
			slog.Debug("Failed to write to WS client, closing connection", "error", err)
			conn.Close()
			delete(s.conns, conn)
		}
	}
}

// handleIndex serves the premium editor and rendering web interface.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		slog.Error("Failed to render index HTML", "error", err)
	}
}

// handleWS handles incoming WebSocket connections, upgrading them and
// running the message loop.
func (s *Server) handleWS(ws *websocket.Conn) {
	s.mu.Lock()
	s.conns[ws] = true
	initialCode := s.code
	s.mu.Unlock()

	// Send initial state to the newly connected editor client.
	initMsg := WSMessage{
		Type: "init",
		Code: initialCode,
	}
	payload, err := json.Marshal(initMsg)
	if err == nil {
		_, _ = ws.Write(payload)
	}

	defer func() {
		s.mu.Lock()
		delete(s.conns, ws)
		s.mu.Unlock()
		ws.Close()
	}()

	var buf []byte = make([]byte, 65536)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err != io.EOF {
				slog.Debug("WebSocket read error", "error", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			slog.Error("Failed to parse WebSocket message", "error", err)
			continue
		}

		switch msg.Type {
		case "update":
			s.mu.Lock()
			s.code = msg.Code
			s.mu.Unlock()
			if s.onChange != nil {
				s.onChange(msg.Code)
			}
		case "sync":
			s.mu.Lock()
			s.code = msg.Code
			s.mu.Unlock()
			if s.onChange != nil {
				s.onChange(msg.Code)
			}
			if s.onSync != nil {
				s.onSync(msg.Code)
			}
		}
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Casspr Mermaid Live Editor</title>
    <style>
        :root {
            --bg-base: #090d16;
            --bg-glass: rgba(17, 24, 39, 0.7);
            --fg-base: #f3f4f6;
            --fg-muted: #9ca3af;
            --primary: #6366f1;
            --primary-hover: #4f46e5;
            --border-glass: rgba(255, 255, 255, 0.08);
            --font-sans: 'Inter', system-ui, sans-serif;
            --font-mono: 'Fira Code', 'Courier New', monospace;
            --shadow-glow: 0 0 20px rgba(99, 102, 241, 0.2);
        }
        
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            background: linear-gradient(135deg, #090d16, #111827);
            color: var(--fg-base);
            font-family: var(--font-sans);
            height: 100vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
        
        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem 2rem;
            background: var(--bg-glass);
            backdrop-filter: blur(12px);
            border-bottom: 1px solid var(--border-glass);
            z-index: 10;
        }
        
        .logo {
            font-size: 1.25rem;
            font-weight: 700;
            background: linear-gradient(to right, #818cf8, #a78bfa);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .actions {
            display: flex;
            gap: 1rem;
            align-items: center;
        }
        
        .btn {
            background: var(--primary);
            color: white;
            border: none;
            padding: 0.5rem 1.25rem;
            border-radius: 6px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .btn:hover {
            background: var(--primary-hover);
            box-shadow: var(--shadow-glow);
            transform: translateY(-1px);
        }
        
        .btn:active {
            transform: translateY(0);
        }
        
        .btn-secondary {
            background: rgba(255, 255, 255, 0.08);
            border: 1px solid var(--border-glass);
        }
        
        .btn-secondary:hover {
            background: rgba(255, 255, 255, 0.15);
            box-shadow: none;
        }
        
        .status {
            font-size: 0.875rem;
            color: var(--fg-muted);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        
        .status-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: #10b981;
            box-shadow: 0 0 8px #10b981;
        }
        
        .status-dot.disconnected {
            background: #ef4444;
            box-shadow: 0 0 8px #ef4444;
        }
        
        main {
            flex: 1;
            display: flex;
            overflow: hidden;
        }
        
        .panel {
            display: flex;
            flex-direction: column;
            height: 100%;
        }
        
        .editor-panel {
            width: 40%;
            border-right: 1px solid var(--border-glass);
            background: rgba(15, 23, 42, 0.3);
        }
        
        .preview-panel {
            width: 60%;
            background: #0d1117;
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 2rem;
            overflow: auto;
        }
        
        .panel-header {
            padding: 0.75rem 1.5rem;
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--fg-muted);
            border-bottom: 1px solid var(--border-glass);
            background: rgba(0, 0, 0, 0.1);
        }
        
        #editor {
            flex: 1;
            width: 100%;
            height: 100%;
            background: transparent;
            color: #e5e7eb;
            font-family: var(--font-mono);
            font-size: 0.95rem;
            padding: 1.5rem;
            border: none;
            resize: none;
            outline: none;
            line-height: 1.6;
        }
        
        #preview {
            max-width: 100%;
            max-height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        
        #error-overlay {
            position: absolute;
            bottom: 1.5rem;
            left: 1.5rem;
            right: 1.5rem;
            background: rgba(220, 38, 38, 0.95);
            backdrop-filter: blur(4px);
            border: 1px solid rgba(248, 113, 113, 0.2);
            padding: 1rem;
            border-radius: 8px;
            color: #fef2f2;
            font-size: 0.875rem;
            font-family: var(--font-mono);
            opacity: 0;
            transform: translateY(10px);
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            pointer-events: none;
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
        }
        
        #error-overlay.visible {
            opacity: 1;
            transform: translateY(0);
        }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
</head>
<body>
    <header>
        <div class="logo">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 3v18h18"/>
                <path d="m19 9-5 5-4-4-3 3"/>
            </svg>
            Casspr Mermaid Live Editor
        </div>
        <div class="actions">
            <div class="status">
                <div id="status-dot" class="status-dot"></div>
                <span id="status-text">Connected</span>
            </div>
            <button class="btn btn-secondary" onclick="copyMarkdown()">Copy Markdown</button>
            <button class="btn" onclick="syncBack(event)">Sync back to Casspr</button>
        </div>
    </header>
    
    <main>
        <div class="panel editor-panel">
            <div class="panel-header">Mermaid Source</div>
            <textarea id="editor" spellcheck="false" placeholder="Enter mermaid code here..."></textarea>
        </div>
        
        <div class="panel preview-panel">
            <div class="panel-header">Visual Preview</div>
            <div id="preview"></div>
            <div id="error-overlay"></div>
        </div>
    </main>

    <script>
        mermaid.initialize({
            startOnLoad: false,
            theme: 'dark',
            securityLevel: 'loose'
        });

        const editor = document.getElementById('editor');
        const preview = document.getElementById('preview');
        const errorOverlay = document.getElementById('error-overlay');
        const statusDot = document.getElementById('status-dot');
        const statusText = document.getElementById('status-text');

        let ws;
        let renderTimeout;

        function connectWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            ws = new WebSocket(protocol + '//' + window.location.host + '/ws');
            
            ws.onopen = () => {
                statusDot.className = 'status-dot';
                statusText.innerText = 'Connected';
            };
            
            ws.onclose = () => {
                statusDot.className = 'status-dot disconnected';
                statusText.innerText = 'Disconnected - Reconnecting...';
                setTimeout(connectWebSocket, 2000);
            };
            
            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                if (data.type === 'init' || data.type === 'update') {
                    editor.value = data.code;
                    renderDiagram();
                }
            };
        }

        async function renderDiagram() {
            const code = editor.value.trim();
            if (!code) {
                preview.innerHTML = '';
                errorOverlay.classList.remove('visible');
                return;
            }
            
            try {
                errorOverlay.classList.remove('visible');
                const { svg } = await mermaid.render('mermaid-svg', code);
                preview.innerHTML = svg;
            } catch (err) {
                errorOverlay.innerText = err.message || err;
                errorOverlay.classList.add('visible');
                
                const errElement = document.getElementById('dmermaid-svg');
                if (errElement) errElement.remove();
            }
        }

        editor.addEventListener('input', () => {
            clearTimeout(renderTimeout);
            renderTimeout = setTimeout(() => {
                renderDiagram();
                sendUpdate(editor.value);
            }, 300);
        });

        function sendUpdate(code) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: 'update',
                    code: code
                }));
            }
        }

        function syncBack(event) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: 'sync',
                    code: editor.value
                }));
                const btn = event.target;
                const oldText = btn.innerText;
                btn.innerText = '✓ Synced!';
                btn.style.background = '#10b981';
                setTimeout(() => {
                    btn.innerText = oldText;
                    btn.style.background = '';
                }, 2000);
            }
        }

        function copyMarkdown() {
            const tick = String.fromCharCode(96);
            const markdown = tick + tick + tick + "mermaid\n" + editor.value + "\n" + tick + tick + tick;
            navigator.clipboard.writeText(markdown).then(() => {
                alert('Copied to clipboard!');
            });
        }

        connectWebSocket();
    </script>
</body>
</html>
`
