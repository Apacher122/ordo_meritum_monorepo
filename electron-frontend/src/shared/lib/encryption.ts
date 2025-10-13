import { apiRequest } from "@/shared/utils/requests";


class PublicKeyService {
  private publicKeyPem: string | null = null;
  private fetchPromise: Promise<void> | null = null;

  private async fetchKey(): Promise<void> {
    try {

      const res = await apiRequest("public-key", {
        method: "GET"
      }) as any;
      this.publicKeyPem = res.publicKey;

      console.log("Fetched public key.");
      console.log(this.publicKeyPem);
    } catch (err) {
      console.error("Failed to fetch public key", err);
      throw new Error("Could not load encryption key, please try again later.");
    }
  }

  public async getPublicKey(): Promise<string> {
    if (this.publicKeyPem) {
      return this.publicKeyPem;
    }
    if (!this.fetchPromise) {
      this.fetchPromise = this.fetchKey();
    }
    await this.fetchPromise;
    if (!this.publicKeyPem) {
      throw new Error("Public key is not available after fetch.");
    }
    return this.publicKeyPem;
  }
}

const keyService = new PublicKeyService();

function pemToArrayBuffer(pem: string): ArrayBuffer {
  const b64 = pem
    .replace(/-----BEGIN RSA PUBLIC KEY-----/, "")
    .replace(/-----END RSA PUBLIC KEY-----/, "")
    .replace(/\s/g, "");
  const binary = atob(b64);
  const buf = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    buf[i] = binary.charCodeAt(i);
  }
  return buf.buffer;
}

function arrayBufferToBase64(buffer: ArrayBuffer): string {
  let binary = "";
  const bytes = new Uint8Array(buffer);
  const len = bytes.byteLength;
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

/**
 * Encrypts a string of data using the fetched public key.
 * @param {string} data - The string to encrypt (e.g., the API key).
 * @returns {Promise<string>} A promise that resolves to the Base64 encoded encrypted string.
 */
export async function encryptData(data: string): Promise<string> {
  const pem = await keyService.getPublicKey();
  const key = await crypto.subtle.importKey(
    "spki",
    pemToArrayBuffer(pem),
    { name: "RSA-OAEP", hash: "SHA-256" },
    false,
    ["encrypt"]
  );
  const encoded = new TextEncoder().encode(data);
  const encrypted = await crypto.subtle.encrypt(
    { name: "RSA-OAEP" },
    key,
    encoded
  );
  return arrayBufferToBase64(encrypted);
}