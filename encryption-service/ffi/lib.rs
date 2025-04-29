use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;

use crate::EncryptionService;

#[no_mangle]
pub extern "C" fn encryption_service_new() -> *mut EncryptionService {
    let service = EncryptionService::new();
    Box::into_raw(Box::new(service))
}

#[no_mangle]
pub extern "C" fn encryption_service_free(service: *mut EncryptionService) {
    if !service.is_null() {
        unsafe {
            Box::from_raw(service);
        }
    }
}

#[no_mangle]
pub extern "C" fn encryption_service_encrypt(service: *mut EncryptionService, plaintext: *const c_char) -> *mut c_char {
    if service.is_null() || plaintext.is_null() {
        return ptr::null_mut();
    }

    let service = unsafe { &*service };
    let plaintext = unsafe { CStr::from_ptr(plaintext) }.to_bytes();
    let ciphertext = service.encrypt(plaintext);

    let ciphertext_cstring = CString::new(ciphertext).unwrap();
    ciphertext_cstring.into_raw()
}

#[no_mangle]
pub extern "C" fn encryption_service_decrypt(service: *mut EncryptionService, ciphertext: *const c_char) -> *mut c_char {
    if service.is_null() || ciphertext.is_null() {
        return ptr::null_mut();
    }

    let service = unsafe { &*service };
    let ciphertext = unsafe { CStr::from_ptr(ciphertext) }.to_bytes();
    let plaintext = service.decrypt(ciphertext);

    let plaintext_cstring = CString::new(plaintext).unwrap();
    plaintext_cstring.into_raw()
}

#[no_mangle]
pub extern "C" fn encryption_service_derive_key(service: *mut EncryptionService, salt: *const c_char, info: *const c_char) -> *mut c_char {
    if service.is_null() || salt.is_null() || info.is_null() {
        return ptr::null_mut();
    }

    let service = unsafe { &*service };
    let salt = unsafe { CStr::from_ptr(salt) }.to_bytes();
    let info = unsafe { CStr::from_ptr(info) }.to_bytes();
    let derived_key = service.derive_key(salt, info);

    let derived_key_cstring = CString::new(derived_key).unwrap();
    derived_key_cstring.into_raw()
}
