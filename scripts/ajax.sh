#!/bin/sh
TIMESTAMP=$(($(date +%s%N)/1000000))
URL="https://www.rusprofile.ru/ajax.php?query=7736207543&action=search&cacheKey=0.$TIMESTAMP"

curl "$URL" \
  -H 'authority: www.rusprofile.ru' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'accept-language: ru-RU,ru;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6' \
  -H 'cache-control: no-cache' \
  -H 'cookie: fbb_s=1; fbb_u=1681175315; _ym_uid=1681175316650019099; _ym_d=1681175316; _ym_isad=2; _ym_visorc=b; _sp_ses.6279=*; _sp_id.6279=6dc9a758-5268-47d1-8a07-59f6724172fc.1681175317.1.1681175357..ff7075fe-c231-44e6-a0b3-5b335b495e3e..dc5b1418-434f-4632-a86e-d85b343481ce.1681175316835.3' \
  -H 'pragma: no-cache' \
  -H 'referer: https://www.rusprofile.ru/' \
  -H 'sec-ch-ua: "Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36' \
  --compressed