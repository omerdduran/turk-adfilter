import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const fileName = searchParams.get('file');

    if (!fileName) {
      return NextResponse.json(
        { error: 'Dosya adı belirtilmedi' },
        { status: 400 }
      );
    }

    // Get the project root (go up from frontend/app/api/backups/download/)
    const projectRoot = path.resolve(process.cwd(), '../');
    const backupDir = path.join(projectRoot, 'yedekler');
    const filePath = path.join(backupDir, fileName);

    // Security check: ensure the file is within the backup directory
    if (!filePath.startsWith(backupDir)) {
      return NextResponse.json(
        { error: 'Geçersiz dosya yolu' },
        { status: 403 }
      );
    }

    // Check if file exists
    if (!fs.existsSync(filePath)) {
      return NextResponse.json(
        { error: 'Dosya bulunamadı' },
        { status: 404 }
      );
    }

    // Read file
    const fileBuffer = fs.readFileSync(filePath);
    
    // Determine content type based on file extension
    let contentType = 'application/octet-stream';
    if (fileName.endsWith('.txt')) {
      contentType = 'text/plain; charset=utf-8';
    }

    // Return file with appropriate headers
    return new NextResponse(fileBuffer, {
      headers: {
        'Content-Type': contentType,
        'Content-Disposition': `attachment; filename="${fileName}"`,
        'Content-Length': fileBuffer.length.toString(),
      },
    });

  } catch (error) {
    console.error('Error downloading backup file:', error);
    return NextResponse.json(
      { 
        error: 'Dosya indirilirken hata oluştu',
        message: error instanceof Error ? error.message : 'Unknown error'
      },
      { status: 500 }
    );
  }
} 