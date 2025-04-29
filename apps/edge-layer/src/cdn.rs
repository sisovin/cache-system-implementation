use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use hyper::{Body, Request, Response};
use hyper::header::{ETag, IfNoneMatch};
use flate2::write::GzEncoder;
use flate2::Compression;
use std::io::prelude::*;

type SharedCache = Arc<Mutex<HashMap<String, (String, String)>>>;

pub struct Cdn {
    cache: SharedCache,
}

impl Cdn {
    pub fn new() -> Self {
        Cdn {
            cache: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    pub async fn handle_request(&self, req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
        let path = req.uri().path().to_string();
        let mut cache = self.cache.lock().unwrap();

        if let Some((etag, content)) = cache.get(&path) {
            if let Some(if_none_match) = req.headers().get(IfNoneMatch::name()) {
                if if_none_match.to_str().unwrap() == etag {
                    return Ok(Response::builder()
                        .status(304)
                        .body(Body::empty())
                        .unwrap());
                }
            }

            let mut encoder = GzEncoder::new(Vec::new(), Compression::default());
            encoder.write_all(content.as_bytes()).unwrap();
            let compressed_content = encoder.finish().unwrap();

            return Ok(Response::builder()
                .header(ETag::name(), etag)
                .header("Content-Encoding", "gzip")
                .body(Body::from(compressed_content))
                .unwrap());
        }

        let response = format!("Hello, you requested: {}", path);
        let etag = format!("\"{}\"", md5::compute(&response));

        cache.insert(path.clone(), (etag.clone(), response.clone()));

        let mut encoder = GzEncoder::new(Vec::new(), Compression::default());
        encoder.write_all(response.as_bytes()).unwrap();
        let compressed_content = encoder.finish().unwrap();

        Ok(Response::builder()
            .header(ETag::name(), etag)
            .header("Content-Encoding", "gzip")
            .body(Body::from(compressed_content))
            .unwrap())
    }
}
