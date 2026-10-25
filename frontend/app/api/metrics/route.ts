import { NextResponse } from 'next/server';

// Cache data structure
interface CachedData {
  metrics: Metrics;
  timestamp: number;
}

// Metrics interface
interface Metrics {
  totalDomains: number;
  filterListDomains: number;
  hostsListDomains: number;
  lastUpdated: string;
}

// Cache duration in milliseconds (15 minutes)
const CACHE_DURATION = 15 * 60 * 1000;

// In-memory cache
let metricsCache: CachedData | null = null;

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

// Get Turkey local time
function getTurkeyTime() {
  return new Date().toLocaleString('tr-TR', { 
    timeZone: 'Europe/Istanbul',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
}

export async function GET() {
  // Check if we have cached data that's not expired
  const now = Date.now();
  if (metricsCache && (now - metricsCache.timestamp) < CACHE_DURATION) {
    return NextResponse.json(metricsCache.metrics, {
      headers: {
        'Cache-Control': 'public, max-age=900, stale-while-revalidate=60'
      }
    });
  }

  const filterListUrl = "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt";
  const hostsListUrl = "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/hosts.txt";

  try {
    const [filterList, hostsList] = await Promise.all([
      fetchFilterList(filterListUrl),
      fetchFilterList(hostsListUrl)
    ]);

    const metrics = {
      totalDomains: countDomains(filterList) + countDomains(hostsList),
      filterListDomains: countDomains(filterList),
      hostsListDomains: countDomains(hostsList),
      lastUpdated: getTurkeyTime()
    };

    // Update cache
    metricsCache = {
      metrics,
      timestamp: now
    };

    return NextResponse.json(metrics, {
      headers: {
        'Cache-Control': 'public, max-age=900, stale-while-revalidate=60'
      }
    });
  } catch (error) {
    console.error('Error getting metrics:', error);
    return NextResponse.json(
      { error: 'Failed to fetch metrics' },
      { status: 500 }
    );
  }
} 