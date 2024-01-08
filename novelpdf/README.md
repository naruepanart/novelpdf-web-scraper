# novelpdf

## Copy urls

```js
const rr = $$('#manga-chapters-holder > div.page-content-listing.single-page > div > ul > li > a').map(x => x.getAttribute("href"))
copy(rr)
```