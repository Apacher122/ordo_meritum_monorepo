import { apiRequest } from "@/shared/utils/requests";

/**
 * @class PublicKeyService
 * A singleton service to fetch and cache the server's public RSA key.
 */
class PublicKeyService {
  private publicKeyPem: string | null = null;
  private fetchPromise: Promise<void> | null = null;

  /**
   * Fetches the server's public RSA key.
   * It throws an error if the key could not be loaded.
   * @returns {Promise<void>} A promise that resolves when the key has been loaded.
   */
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

  /**
   * Retrieves the public RSA key in PEM format.
   * It fetches the key on the first call and caches it for subsequent requests.
   * @returns {Promise<string>} A promise that resolves to the public key in PEM format.
   * @throws {Error} If the key is unavailable after the fetch attempt.
   */
  public async getPublicKey(): Promise<string> {
    if (this.publicKeyPem) {
      return this.publicKeyPem;
    }
    this.fetchPromise ??= this.fetchKey();
    await this.fetchPromise;
    if (!this.publicKeyPem) {
      throw new Error("Public key is not available after fetch.");
    }
    return this.publicKeyPem;
  }
}

const keyService = new PublicKeyService();

/**
 * Converts a PEM-formatted public key string into an ArrayBuffer.
 * @param {string} pem The PEM string.
 * @returns {ArrayBuffer} The key as an ArrayBuffer.
 */
function pemToArrayBuffer(pem: string): ArrayBuffer {
  const b64 = pem
    .replace(/-----BEGIN RSA PUBLIC KEY-----/, "")
    .replace(/-----END RSA PUBLIC KEY-----/, "")
    .replaceAll(/\s/g, "");
  const binary = atob(b64);
  const buf = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    buf[i] = binary.codePointAt(i)!;
  }
  return buf.buffer;
}

/**
 * Converts an ArrayBuffer to a Base64 encoded string.
 * @param {ArrayBuffer} buffer The ArrayBuffer to convert.
 * @returns {string} The Base64 encoded string.
 */
function arrayBufferToBase64(buffer: ArrayBuffer): string {
  let binary = "";
  const bytes = new Uint8Array(buffer);
  const len = bytes.byteLength;
  for (let i = 0; i < len; i++) {
    binary += String.fromCodePoint(bytes[i]);
  }
  return btoa(binary);
}

/**
 * Encrypts a string using the server's public RSA key with RSA-OAEP padding.
 * @param {string} data The plain text string to encrypt.
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