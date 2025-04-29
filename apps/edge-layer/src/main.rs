use hyper::service::{make_service_fn, service_fn};
use hyper::{Body, Request, Response, Server};
use lru::LruCache;
use prometheus::{Encoder, TextEncoder, register_counter, register_histogram, Counter, Histogram};
use std::sync::{Arc, Mutex};
use tokio::runtime::Runtime;

type SharedCache = Arc<Mutex<LruCache<String, String>>>;

async fn handle_request(req: Request<Body>, cache: SharedCache, request_counter: Counter, request_histogram: Histogram) -> Result<Response<Body>, hyper::Error> {
    let path = req.uri().path().to_string();

    // Start timer for request duration
    let timer = request_histogram.start_timer();

    // Increment request counter
    request_counter.inc();

    // Check cache
    let mut cache = cache.lock().unwrap();
    if let Some(response) = cache.get(&path) {
        return Ok(Response::new(Body::from(response.clone())));
    }

    // Generate response
    let response = format!("Hello, you requested: {}", path);

    // Store response in cache
    cache.put(path.clone(), response.clone());

    // Stop timer for request duration
    timer.observe_duration();

    Ok(Response::new(Body::from(response)))
}

async fn serve_metrics() -> Result<Response<Body>, hyper::Error> {
    let encoder = TextEncoder::new();
    let metric_families = prometheus::gather();
    let mut buffer = Vec::new();
    encoder.encode(&metric_families, &mut buffer).unwrap();
    Ok(Response::new(Body::from(buffer)))
}

fn main() {
    // Initialize LRU cache
    let cache = Arc::new(Mutex::new(LruCache::new(100)));

    // Initialize Prometheus metrics
    let request_counter = register_counter!("requests_total", "Total number of requests").unwrap();
    let request_histogram = register_histogram!("request_duration_seconds", "Request duration in seconds").unwrap();

    // Create the runtime
    let rt = Runtime::new().unwrap();

    rt.block_on(async {
        // Create the service
        let make_svc = make_service_fn(|_conn| {
            let cache = cache.clone();
            let request_counter = request_counter.clone();
            let request_histogram = request_histogram.clone();
            async {
                Ok::<_, hyper::Error>(service_fn(move |req| {
                    if req.uri().path() == "/metrics" {
                        serve_metrics()
                    } else {
                        handle_request(req, cache.clone(), request_counter.clone(), request_histogram.clone())
                    }
                }))
            }
        });

        // Create the server
        let addr = ([127, 0, 0, 1], 3000).into();
        let server = Server::bind(&addr).serve(make_svc);

        // Run the server
        if let Err(e) = server.await {
            eprintln!("server error: {}", e);
        }
    });
}
