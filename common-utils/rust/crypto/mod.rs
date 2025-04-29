use aes_gcm::aead::{Aead, KeyInit, OsRng};
use aes_gcm::{Aes256Gcm, Key, Nonce}; // Or `Aes128Gcm`
use hkdf::Hkdf;
use sha2::Sha256;
use criterion::{criterion_group, criterion_main, Criterion};

pub struct CryptoService {
    key: Key<Aes256Gcm>,
}

impl CryptoService {
    pub fn new() -> Self {
        let mut key = [0u8; 32];
        OsRng.fill_bytes(&mut key);
        let key = Key::<Aes256Gcm>::from_slice(&key);
        CryptoService {
            key: key.clone(),
        }
    }

    pub fn encrypt(&self, plaintext: &[u8]) -> Vec<u8> {
        let cipher = Aes256Gcm::new(&self.key);
        let nonce = Nonce::from_slice(b"unique nonce"); // 96-bits; unique per message
        cipher.encrypt(nonce, plaintext).expect("encryption failure!")
    }

    pub fn decrypt(&self, ciphertext: &[u8]) -> Vec<u8> {
        let cipher = Aes256Gcm::new(&self.key);
        let nonce = Nonce::from_slice(b"unique nonce"); // 96-bits; unique per message
        cipher.decrypt(nonce, ciphertext).expect("decryption failure!")
    }

    pub fn derive_key(&self, salt: &[u8], info: &[u8]) -> Vec<u8> {
        let hk = Hkdf::<Sha256>::new(Some(salt), &self.key);
        let mut okm = [0u8; 32];
        hk.expand(info, &mut okm).expect("HKDF expand failed");
        okm.to_vec()
    }
}

fn benchmark_encrypt(c: &mut Criterion) {
    let service = CryptoService::new();
    let plaintext = b"Benchmarking encryption performance";

    c.bench_function("encrypt", |b| {
        b.iter(|| {
            service.encrypt(plaintext);
        })
    });
}

fn benchmark_decrypt(c: &mut Criterion) {
    let service = CryptoService::new();
    let ciphertext = service.encrypt(b"Benchmarking decryption performance");

    c.bench_function("decrypt", |b| {
        b.iter(|| {
            service.decrypt(&ciphertext);
        })
    });
}

criterion_group!(benches, benchmark_encrypt, benchmark_decrypt);
criterion_main!(benches);
