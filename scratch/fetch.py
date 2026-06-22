import json

file_path = "/Users/ibranraeen/.gemini/antigravity-ide/brain/80ffef8b-6326-40e7-b694-7c6abc3df141/.system_generated/steps/326/content.md"

try:
    with open(file_path, "r") as f:
        content = f.read()
    
    # Strip headers if present
    json_start = content.find("{\"total_count\"")
    if json_start != -1:
        json_data = json.loads(content[json_start:])
        print("Total runs found:", len(json_data.get("workflow_runs", [])))
        for run in json_data.get("workflow_runs", []):
            if "release.yml" in run.get("path", ""):
                print(f"ID: {run['id']}, Name: {run['name']}, Event: {run['event']}, Status: {run['status']}, Conclusion: {run['conclusion']}, SHA: {run['head_sha']}, Display Title: {run['display_title']}")
    else:
        print("Could not find start of JSON in the file.")
except Exception as e:
    print("Error:", e)
