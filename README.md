# Golang opentracing example with gin (WORKING!!)  

--- 

> ### WANT TO DO  


```  
             [Span A (Service1)]  
                     |  
         +-----------+-----------+  
         |                       |  
 [Span B (Service2)]      [Span C (Service 3)] >>> [Span D (Service 4)]  
         |                   
 [Span E (Service5)]
```



---


this project is example of opentracing-go with gin, ...  

---  

## Overview

Suppose that we have 5 services like below tree.

> Service to service calls

```  
             [Span A (Service1)]  
                     |  
         +-----------+-----------+  
         |                       |  
 [Span B (Service2)]      [Span C (Service 3)] >>> [Span D (Service 4)]  
         |                   
 [Span E (Service5)]
```  

## Getting started  

> Run with docker compose  

```bash
$ docker-compose up

// if u want build image again, use --build tag
$ docker-compose up --build

// jaeger ui
// http://localhost:16686/
```
