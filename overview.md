# current overview

this is just so i can see what i did.

## domain

ajitesh.cc domain -> bought on namecheap with NS for cloudflare dns.

## dns records related to services

ajitesh.cc -> four A records to github pages ips

www.ajitesh.cc -> CNAME to `tiredkangaroo.github.io` (github pages)

server.ajitesh.cc -> cloudflare proxied A record to my home ip addr.

photos.ajitesh.cc -> managed by cloudflare r2.

potentially adding AAAA records for github pages and home ip later.

## frontend

written in sveltekit. deployed using github pages and github actions.

## backend

written in go. sql go generated code with sqlc. a container on a raspi at home.

## database

psql. it's a container on the same raspi as the backend.

## object storage

cloudflare r2. used for storing photos for photography section.

## services that may require payment

- namecheap for domain (currently paying, requires renewal by oct 5, 2026)
- cloudflare r2 for object storage (potentially paying, free tier should be enough)
