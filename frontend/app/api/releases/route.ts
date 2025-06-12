import { NextResponse } from "next/server";

const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
const REPO_OWNER = "omerdduran";
const REPO_NAME = "turk-adfilter";

// Cache for releases
let releasesCache: {
  data: any[];
  timestamp: number;
} | null = null;

const CACHE_DURATION = 30 * 60 * 1000; // 30 minutes

export interface GitHubRelease {
  id: number;
  tag_name: string;
  name: string;
  body: string;
  published_at: string;
  html_url: string;
  draft: boolean;
  prerelease: boolean;
  assets: Array<{
    name: string;
    download_count: number;
    size: number;
    browser_download_url: string;
  }>;
}

// Function to clean release body by removing QA section
const cleanReleaseBody = (body: string): string => {
  if (!body) return body;

  // Remove the "Otomatik Kalite Kontrolleri" section and everything after it
  const qaSectionStart = body.indexOf("## 🧪 Otomatik Kalite Kontrolleri");
  if (qaSectionStart !== -1) {
    return body.substring(0, qaSectionStart).trim();
  }

  // Also check for English version
  const qaSectionStartEn = body.indexOf("## 🧪 Automatic Quality Controls");
  if (qaSectionStartEn !== -1) {
    return body.substring(0, qaSectionStartEn).trim();
  }

  return body;
};

export async function GET(request: Request) {
  if (!GITHUB_TOKEN) {
    return NextResponse.json(
      { message: "GitHub token is not configured" },
      { status: 500 }
    );
  }

  const { searchParams } = new URL(request.url);
  const page = parseInt(searchParams.get("page") || "1");
  const perPage = parseInt(searchParams.get("per_page") || "5");
  const all = searchParams.get("all") === "true";

  try {
    // Check cache first
    const now = Date.now();
    if (releasesCache && now - releasesCache.timestamp < CACHE_DURATION) {
      const cachedReleases = releasesCache.data;

      if (all) {
        return NextResponse.json(cachedReleases, {
          headers: {
            "Cache-Control": "public, max-age=1800, stale-while-revalidate=300",
          },
        });
      }

      const startIndex = (page - 1) * perPage;
      const endIndex = startIndex + perPage;
      const paginatedReleases = cachedReleases.slice(startIndex, endIndex);

      return NextResponse.json(
        {
          releases: paginatedReleases,
          total: cachedReleases.length,
          page,
          perPage,
          hasMore: endIndex < cachedReleases.length,
        },
        {
          headers: {
            "Cache-Control": "public, max-age=1800, stale-while-revalidate=300",
          },
        }
      );
    }

    // Fetch from GitHub if cache is expired or doesn't exist
    const response = await fetch(
      `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=100&page=1`,
      {
        headers: {
          Authorization: `Bearer ${GITHUB_TOKEN}`,
          Accept: "application/vnd.github.v3+json",
        },
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Failed to fetch releases");
    }

    const releases: GitHubRelease[] = await response.json();

    // Filter out drafts, sort by published date, and clean body
    const publishedReleases = releases
      .filter((release) => !release.draft)
      .map((release) => ({
        ...release,
        body: cleanReleaseBody(release.body),
      }))
      .sort(
        (a, b) =>
          new Date(b.published_at).getTime() -
          new Date(a.published_at).getTime()
      );

    // Update cache
    releasesCache = {
      data: publishedReleases,
      timestamp: now,
    };

    if (all) {
      return NextResponse.json(publishedReleases, {
        headers: {
          "Cache-Control": "public, max-age=1800, stale-while-revalidate=300",
        },
      });
    }

    const startIndex = (page - 1) * perPage;
    const endIndex = startIndex + perPage;
    const paginatedReleases = publishedReleases.slice(startIndex, endIndex);

    return NextResponse.json(
      {
        releases: paginatedReleases,
        total: publishedReleases.length,
        page,
        perPage,
        hasMore: endIndex < publishedReleases.length,
      },
      {
        headers: {
          "Cache-Control": "public, max-age=1800, stale-while-revalidate=300",
        },
      }
    );
  } catch (error) {
    console.error("Error fetching releases:", error);
    return NextResponse.json(
      {
        message:
          error instanceof Error ? error.message : "Failed to fetch releases",
      },
      { status: 500 }
    );
  }
}
