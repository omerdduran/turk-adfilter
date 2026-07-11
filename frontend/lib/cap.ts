// Cap (self-hosted CAPTCHA) yapılandırması.
//
// SITE KEY PUBLIC'tir — widget'ta görünür, burada tutulması güvenlidir.
// SECRET KEY burada ASLA bulunmaz; yalnızca sunucu ortamında `CAP_SECRET`
// olarak okunur (frontend/app/api/issues/route.ts). Git'e secret koyulmaz.

export const CAP_INSTANCE = "https://captcha.reklamsiz-turkiye.com";
export const CAP_SITE_KEY = "b9939d9202";

// Widget'ın challenge istediği endpoint (istemci tarafı).
export const CAP_WIDGET_ENDPOINT = `${CAP_INSTANCE}/${CAP_SITE_KEY}/`;

// Sunucunun token doğruladığı endpoint (sunucu tarafı — secret ile).
export const CAP_SITEVERIFY_URL = `${CAP_INSTANCE}/${CAP_SITE_KEY}/siteverify`;
