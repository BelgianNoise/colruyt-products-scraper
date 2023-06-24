# colruyt-products-scraper

An application written in Go that scrapes Colruyt's API to retrieve all product listings.

Data is scraper every night using a cronjob running in Github Actions and uploaded to a public Google Cloud Storage bucket, named `colruyt-products` and hosted in `us-east1`, for everyone to use (and for easier acccess for me later).
To overcome rate limiting and other surprisingly well implemented anti-scrapeing measures I am using publicly available online proxies. (only because the Colruyt API has some weird behaviour)

I wanted to create this because I find it very interesting to know which products have risen by how much in price, and to know when to possibly look for a cheaper alternative.

All data is also stored in a PostgreSQL database. This one is however not publicly available (for now) because of possible cost implications.

### The end goal is:
- to create an accompanying frontend that
  - allows for filtering and viewing product listings from any date.
  - lists price changes
  - compare prices from any date available
  - (Need to fill the bucket with some data first)
- create some kind of notifications system to notify of the biggest price changes of that day or week (maybe useing [Resend](https://resend.com/) ? )
