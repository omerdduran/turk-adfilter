import { NextResponse } from "next/server";

export async function GET() {
  try {
    // Use the existing releases API to get all releases
    const response = await fetch(
      `${
        process.env.NEXTAUTH_URL || "http://localhost:3000"
      }/api/releases?all=true`
    );

    if (!response.ok) {
      throw new Error("Failed to fetch releases");
    }

    const releases = await response.json();
    const rssContent = generateRSS(releases);

    return new NextResponse(rssContent, {
      headers: {
        "Content-Type": "text/xml; charset=utf-8",
        "Cache-Control": "public, max-age=1800, stale-while-revalidate=300",
      },
    });
  } catch (error) {
    console.error("Error generating RSS feed:", error);
    return NextResponse.json(
      {
        message:
          error instanceof Error
            ? error.message
            : "Failed to generate RSS feed",
      },
      { status: 500 }
    );
  }
}

function generateRSS(releases: any[]): string {
  const baseUrl = "https://reklamsiz-turkiye.com";
  const currentDate = new Date().toISOString();

  const rssItems = releases
    .map((release) => {
      const pubDate = new Date(release.published_at).toUTCString();
      const title = `${release.tag_name} - ${release.name || release.tag_name}`;

      return `
      <item>
        <title>${escapeXml(title)}</title>
        <link>${release.html_url}</link>
        <guid>${release.html_url}</guid>
        <pubDate>${pubDate}</pubDate>
      </item>
    `;
    })
    .join("");

  return `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Turk-AdFilter Releases</title>
    <link>${baseUrl}/releases</link>
    <atom:link href="${baseUrl}/api/releases/rss" rel="self" type="application/rss+xml" />
    <description>Türkiye merkezli reklam ve izleyici engelleme listesi Turk-AdFilter'ın en son sürümleri</description>
    <language>tr</language>
    <lastBuildDate>${currentDate}</lastBuildDate>
    <generator>Turk-AdFilter RSS Generator</generator>
    ${rssItems}
  </channel>
</rss>`;
}

function escapeXml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;");
}
