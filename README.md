# novelpdf-web-scraper

 ```js
const selectedLinks = document.querySelectorAll('#manga-chapters-holder > div.page-content-listing.single-page > div > ul > li > a');
const urls = Array.from(selectedLinks).map(link => link.href);
copy(urls);
```