# Bugs known

1. asyncio http server will fail after some unspecified time (Reproducible)
   1. Seems to be related to the number of async connections that the http server can handle at a certain time