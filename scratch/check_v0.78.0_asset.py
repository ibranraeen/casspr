import urllib.request
import urllib.error

url = "https://github.com/ibranraeen/casspr/releases/download/v0.78.0/casspr-windows-amd64.exe"

req = urllib.request.Request(url, method="HEAD")
try:
    with urllib.request.urlopen(req) as resp:
        print("Status:", resp.status)
        print("Headers:")
        for k, v in resp.getheaders():
            print(f"  {k}: {v}")
except urllib.error.HTTPError as e:
    print("HTTP Error:", e.code)
except Exception as e:
    print("Error:", e)
