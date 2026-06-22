import json

filepath = "/Users/ibranraeen/.gemini/antigravity-ide/brain/80ffef8b-6326-40e7-b694-7c6abc3df141/.system_generated/steps/86/content.md"

with open(filepath, "r") as f:
    lines = f.readlines()

# Find JSON start (after ---)
json_start = 0
for idx, line in enumerate(lines):
    if line.strip() == "---":
        json_start = idx + 1
        break

json_str = "".join(lines[json_start:]).strip()
data = json.loads(json_str)

if isinstance(data, list):
    # This might be a list of releases
    release = data[0]
else:
    release = data

print("Release:", release.get("tag_name"))
print("Assets:")
for asset in release.get("assets", []):
    print(f"- {asset.get('name')}: {asset.get('browser_download_url')}")
