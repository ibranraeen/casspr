import fs from 'fs';

async function fetchLogs() {
  const url = 'https://api.github.com/repos/ibranraeen/casspr/actions/jobs/82682109108/logs';
  console.log('Fetching logs from:', url);
  const res = await fetch(url, {
    headers: {
      'Authorization': `Bearer ${process.env.GITHUB_TOKEN}`,
      'Accept': 'application/vnd.github.v3+json',
      'User-Agent': 'Node.js'
    }
  });
  console.log('Status:', res.status);
  if (res.status !== 200) {
    const errText = await res.text();
    console.error('Error body:', errText);
    return;
  }
  const text = await res.text();
  // Write to a local file so we can view it
  fs.writeFileSync('scratch/job_logs.txt', text);
  console.log('Successfully wrote logs to scratch/job_logs.txt. Length:', text.length);
}

fetchLogs().catch(console.error);
