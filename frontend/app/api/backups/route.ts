import { NextResponse } from 'next/server';

interface BackupInfo {
  backup_file: string;
  timestamp: string;
  size: number;
  hash: string;
  created_at: string;
}

interface FileBackups {
  last_hash: string;
  last_backup: string;
  backups: BackupInfo[];
}

interface BackupMetadata {
  [filename: string]: FileBackups;
}

// Cache for backup data
let backupCache: {
  data: any;
  timestamp: number;
} | null = null;

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

export async function GET() {
  try {
    // Check cache first
    const now = Date.now();
    if (backupCache && now - backupCache.timestamp < CACHE_DURATION) {
      return NextResponse.json(backupCache.data, {
        headers: {
          'Cache-Control': 'public, max-age=300, stale-while-revalidate=60'
        }
      });
    }

    // Fetch metadata from GitHub
    const metadataUrl = 'https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/yedekler/backup_metadata.json';
    
    const response = await fetch(metadataUrl);
    
    if (!response.ok) {
      if (response.status === 404) {
        return NextResponse.json({
          error: 'Henüz yedek dosyası oluşturulmamış',
          backups: {},
          stats: {
            totalFiles: 0,
            totalBackups: 0,
            lastBackup: null
          }
        });
      }
      throw new Error(`GitHub'dan metadata alınamadı: ${response.statusText}`);
    }

    const metadataContent = await response.text();
    const metadata: BackupMetadata = JSON.parse(metadataContent);

    // Process backup data
    const processedBackups: any = {};
    let totalBackups = 0;
    let latestTimestamp = '';

    Object.entries(metadata).forEach(([filename, fileData]) => {
      const backups = fileData.backups || [];
      
      // Sort backups by timestamp (newest first)
      const sortedBackups = backups.sort((a, b) => b.timestamp.localeCompare(a.timestamp));
      
      processedBackups[filename] = {
        lastBackup: fileData.last_backup,
        lastHash: fileData.last_hash,
        backups: sortedBackups.map(backup => ({
          ...backup,
          downloadUrl: `https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/yedekler/${backup.backup_file}`,
          githubUrl: `https://github.com/omerdduran/turk-adfilter/blob/main/yedekler/${backup.backup_file}`,
          codebergUrl: `https://codeberg.org/omerdduran/turk-adfilter/src/branch/main/yedekler/${backup.backup_file}`,
          codebergRawUrl: `https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/yedekler/${backup.backup_file}`,
          readableDate: formatTimestamp(backup.timestamp),
          readableSize: formatSize(backup.size)
        }))
      };

      totalBackups += backups.length;
      
      // Find latest timestamp
      if (backups.length > 0) {
        const fileLatest = backups.reduce((latest, backup) => 
          backup.timestamp > latest ? backup.timestamp : latest, '');
        if (fileLatest > latestTimestamp) {
          latestTimestamp = fileLatest;
        }
      }
    });

    const result = {
      backups: processedBackups,
      stats: {
        totalFiles: Object.keys(metadata).length,
        totalBackups,
        lastBackup: latestTimestamp ? formatTimestamp(latestTimestamp) : null,
        lastBackupTimestamp: latestTimestamp
      }
    };

    // Update cache
    backupCache = {
      data: result,
      timestamp: now
    };

    return NextResponse.json(result, {
      headers: {
        'Cache-Control': 'public, max-age=300, stale-while-revalidate=60'
      }
    });

  } catch (error) {
    console.error('Error reading backup data:', error);
    return NextResponse.json(
      { 
        error: 'Yedek verileri okunurken hata oluştu',
        message: error instanceof Error ? error.message : 'Unknown error',
        backups: {},
        stats: {
          totalFiles: 0,
          totalBackups: 0,
          lastBackup: null
        }
      },
      { status: 500 }
    );
  }
}

function formatTimestamp(timestamp: string): string {
  try {
    // Parse timestamp format: YYYYMMDD_HHMMSS
    const year = parseInt(timestamp.substr(0, 4));
    const month = parseInt(timestamp.substr(4, 2)) - 1; // JS months are 0-indexed
    const day = parseInt(timestamp.substr(6, 2));
    const hour = parseInt(timestamp.substr(9, 2));
    const minute = parseInt(timestamp.substr(11, 2));
    const second = parseInt(timestamp.substr(13, 2));

    const date = new Date(year, month, day, hour, minute, second);
    
    return date.toLocaleString('tr-TR', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  } catch (error) {
    return timestamp; // Fallback to original if parsing fails
  }
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
} 