# colruyt-products-scraper

An application written in Go that scrapes Colruyt's API to retrieve all product listings. 

Data is scraped every night using a cronjob running in Github Actions and uploaded to a public Google Cloud Storage bucket, named `colruyt-products` and hosted in `us-east1`, for everyone to use (and for easier acccess for me later).

I wanted to create this because I find it very interesting to know which products have risen by how much in price, and to know when to possibly look for a cheaper alternative.

All data is also stored in a PostgreSQL database. This one is however not publicly available (for now) because of possible cost implications.

A frontend to display all the data can be found at: https://colruyt-prijzen.nasaj.be/ - source: https://github.com/BelgianNoise/colruyt-price-history

# Anti-Bot protections and how to beat them

Disclaimer: I am no expert on anti bot protection.

The Colruyt API implements some very common anti scraping mechanisms like `rate limiting`, but it also enforces extra rules based on cookies and sessions which I do not fully understand yet.

The rate limiting also depends on whether you are `an unknown user`, `a user with a session` or `an authenticated user`.
- **Unknown user**: Your IP will be put on a block list for x amount of time after about 10-20 direct API requests. Yikes.
- **user with a session**: You have increased rate limiting, but after a while you will be served the bot detection page and denied access with the session. (No IP ban)
- **authenticated user**: Besides some more increased rate limiting, I did not see any big difference compared to the `user with a session`.

Moral of the story: You are being rate limited to a speed that would not interfere with regular internet browsing. In my opinion, this is way too slow for me to even consider this a viable option for scraping.

Here are a couple ways you can try to circumvent the anti bot measurements:

- **Proxies**: Sending all requests through a large enough pool of rotating proxies is a great solution to the problem. I sed proxies for a long time, after which the Colruyt product API changed and my proxy bill went from $0.50/mo to $5/mo. I didn't want to shell out this money at the time.
   - **Private Proxies**: Use private/paid proxies for the best result. Choice between `datacenter` and `residential` IPs, where `residential` cost way more but is the most effective. (I used [Bright Data](https://brightdata.com/))
   - **Public Proxies**: Some private proxy providers offer a rotating list of free SSL proxies for you to use, you could simply scrape these web pages and use the free proxies as such. (Free proxies are very hit or miss, I moved away from this rather quickly)
- (Currently in use) **Clean sessions**: Spin up a headless browser, navigate to colruyt.be and let their session management and cookies do their thing. After this, use this close the browser and use its state to convince the API you are that browser. When you reach your limit (which can happen rather quickly, and is very unpredictable), just restart teh whole headless browser shenanigans.
- (unverified) [**Amazon API Gateway**](https://aws.amazon.com/api-gateway/): If I understand correctly, this mimmicks a private proxy with datacenter IPs. You can simply send your request to their service and it will send it out through their servers. While it does not really boast or advertise the rotating IPs, API Gateway should still rotate their IPs, unless asked not to.
 
