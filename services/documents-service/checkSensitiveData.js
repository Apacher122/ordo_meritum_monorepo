
const fs = require('fs');
const path = require('path');

const sensitivePatterns = [
  /API_KEY\s*=\s*['"][\w-]+['"]/i,
  /SECRET\s*=\s*['"][\w-]+['"]/i,
  /PASSWORD\s*=\s*['"][\w-]+['"]/i,
  /ACCESS_TOKEN\s*=\s*['"][\w-]+['"]/i,
  /PRIVATE_KEY\s*=\s*['"][\w-]+['"]/i,
  /AWS_SECRET_ACCESS_KEY\s*=\s*['"][\w-]+['"]/i,
  /GCP_PRIVATE_KEY\s*=\s*['"][\w-]+['"]/i,
  /AZURE_SECRET\s*=\s*['"][\w-]+['"]/i,
];

const scanFile = (filePath) => {
  const content = fs.readFileSync(filePath, 'utf-8');
  for (const pattern of sensitivePatterns) {
    if (pattern.test(content)) {
      console.log(`Sensitive data found in file: ${filePath}`);
      return;
    }
  }
};

const scanDirectory = (dirPath) => {
  const files = fs.readdirSync(dirPath);
  for (const file of files) {
    const fullPath = path.join(dirPath, file);
    if (fs.statSync(fullPath).isDirectory()) {
      scanDirectory(fullPath);
    } else {
      scanFile(fullPath);
    }
  }
};

const directoryToScan = path.resolve(__dirname);
scanDirectory(directoryToScan);