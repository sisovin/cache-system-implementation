from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.trustedhost import TrustedHostMiddleware
from fastapi.middleware.httpsredirect import HTTPSRedirectMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.middleware.cache import CacheMiddleware
from fastapi.middleware.ratelimit import RateLimitMiddleware
from fastapi.middleware.graphql import GraphQLMiddleware
from fastapi.middleware.cachekey import CacheKeyMiddleware
from fastapi.middleware.hierarchicalcache import HierarchicalCacheMiddleware

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Add Trusted Host middleware
app.add_middleware(
    TrustedHostMiddleware,
    allowed_hosts=["*"],
)

# Add HTTPS Redirect middleware
app.add_middleware(
    HTTPSRedirectMiddleware,
)

# Add GZip middleware
app.add_middleware(
    GZipMiddleware,
    minimum_size=1000,
)

# Add Cache middleware
app.add_middleware(
    CacheMiddleware,
    cache_time=60,
)

# Add Rate Limit middleware
app.add_middleware(
    RateLimitMiddleware,
    rate_limit=100,
)

# Add GraphQL middleware
app.add_middleware(
    GraphQLMiddleware,
)

# Add Cache Key middleware
app.add_middleware(
    CacheKeyMiddleware,
)

# Add Hierarchical Cache middleware
app.add_middleware(
    HierarchicalCacheMiddleware,
    cache_levels=["RAM", "Redis", "DB"],
)

@app.get("/")
async def read_root():
    return {"message": "Welcome to the API Gateway"}

@app.get("/items/{item_id}")
async def read_item(item_id: int):
    return {"item_id": item_id}
