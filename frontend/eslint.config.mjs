import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

export default defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    rules: {
      // Türkçe metinlerde kesme işareti/quote karakterleri JSX entity'si
      // gerektirmemeli; kural bu projede sürekli yanlış pozitif üretiyor.
      "react/no-unescaped-entities": "off",
      // Logolar public/ altında sabit boyutlu ve width/height verilmiş;
      // next/image'ın üretimde sharp gerektirmesi self-host imajını kırar.
      "@next/next/no-img-element": "off",
    },
  },
  globalIgnores([".next/**", ".source/**", "next-env.d.ts"]),
]);
