import { NextResponse } from "next/server";

const CAPTCHA_BASE_URL = "https://captcha.reklamsiz-turkiye.com";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { action, ...data } = body;

    let url: string;

    if (action === "generate") {
      url = `${CAPTCHA_BASE_URL}/v2/captcha`;
    } else if (action === "verify") {
      url = `${CAPTCHA_BASE_URL}/v2/answer`;
    } else {
      return NextResponse.json(
        { message: "Invalid action" },
        { status: 400 }
      );
    }

    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`CAPTCHA request failed: ${response.statusText}`);
    }

    const result = await response.json();
    return NextResponse.json(result);
  } catch (error) {
    console.error("CAPTCHA API error:", error);
    return NextResponse.json(
      { message: error instanceof Error ? error.message : "CAPTCHA request failed" },
      { status: 500 }
    );
  }
}

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const id = searchParams.get("id");

    if (!id) {
      return NextResponse.json(
        { message: "CAPTCHA ID is required" },
        { status: 400 }
      );
    }

    const imageUrl = `${CAPTCHA_BASE_URL}/v2/media?id=${id}`;
    
    const response = await fetch(imageUrl);

    if (!response.ok) {
      console.error(`Failed to fetch image: ${response.status} ${response.statusText}`);
      throw new Error(`Failed to fetch CAPTCHA image: ${response.statusText}`);
    }

    const contentType = response.headers.get("content-type"); 
    
    const imageBuffer = await response.arrayBuffer();
    
    return new NextResponse(imageBuffer, {
      headers: {
        "Content-Type": contentType || "image/png",
        "Cache-Control": "no-cache, no-store, must-revalidate",
        "Pragma": "no-cache",
        "Expires": "0",
      },
    });
  } catch (error) {
    console.error("CAPTCHA image fetch error:", error);
    return NextResponse.json(
      { message: error instanceof Error ? error.message : "Failed to fetch CAPTCHA image" },
      { status: 500 }
    );
  }
} 