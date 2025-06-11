import { NextResponse } from "next/server";

const GITHUB_TOKEN = process.env.GITHUB_TOKEN;
const REPO_OWNER = "omerdduran";
const REPO_NAME = "turk-adfilter";

export async function POST(request: Request) {
  if (!GITHUB_TOKEN) {
    return NextResponse.json(
      { message: "GitHub token is not configured" },
      { status: 500 }
    );
  }

  try {
    const requestBody = await request.json();
    const { title, body, labels } = requestBody;

    if (!title || !body) {
      return NextResponse.json(
        { message: "Title and body are required" },
        { status: 400 }
      );
    }

    const response = await fetch(
      `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/issues`,
      {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${GITHUB_TOKEN}`,
          "Accept": "application/vnd.github.v3+json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          title,
          body,
          labels: labels || ["user-feedback"],
        }),
      }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Failed to create issue");
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error creating issue:", error);
    return NextResponse.json(
      { message: error instanceof Error ? error.message : "Failed to create issue" },
      { status: 500 }
    );
  }
} 