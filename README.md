# colruyt-products-scraper

An application written in Go that scrapes Colruyt's API to retrieve all product listings.

Data is scraped every night using a cronjob running in Github Actions and uploaded to a public Google Cloud Storage bucket, named `colruyt-products` and hosted in `us-east1`, for everyone to use (and for easier acccess for me later).
To overcome rate limiting and other surprisingly weird API behaviour I am using private and publicly available online proxies.

I wanted to create this because I find it very interesting to know which products have risen by how much in price, and to know when to possibly look for a cheaper alternative.

All data is also stored in a PostgreSQL database. This one is however not publicly available (for now) because of possible cost implications.

A frontend to display all the data can be found at: https://colruyt-prijzen.nasaj.be/ - source: https://github.com/BelgianNoise/colruyt-price-history
