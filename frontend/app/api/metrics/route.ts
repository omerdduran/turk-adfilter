import { NextResponse } from 'next/server';

async function fetchFilterList(url: string) {
  try {
    const response = await fetch(url);
    const text = await response.text();
    return text;
  } catch (error) {
    console.error('Error fetching filter list:', error);
    return '';
  }
}

function countDomains(text: string) {
  // Remove comments and empty lines
  const lines = text.split('\n').filter(line => 
    line.trim() && !line.startsWith('!') && !line.startsWith('#')
  );
  return lines.length;
}

export async function GET() {
  const filterListUrls = {
    github: "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt",
    codeberg: "https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/turk-adfilter.txt"
  };

  const hostsListUrls = {
    github: "https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/hosts.txt",
    codeberg: "https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/hosts.txt"
  };

  try {
    // Fetch both lists from GitHub (as primary source)
    const [filterList, hostsList] = await Promise.all([
      fetchFilterList(filterListUrls.github),
      fetchFilterList(hostsListUrls.github)
    ]);

    const metrics = {
      totalDomains: countDomains(filterList) + countDomains(hostsList),
      filterListDomains: countDomains(filterList),
      hostsListDomains: countDomains(hostsList),
      lastUpdated: new Date().toISOString()
    };

    return NextResponse.json(metrics);
  } catch (error) {
    console.error('Error getting metrics:', error);
    return NextResponse.json(
      { error: 'Failed to fetch metrics' },
      { status: 500 }
    );
  }
} 