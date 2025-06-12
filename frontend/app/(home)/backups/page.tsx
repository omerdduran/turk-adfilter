"use client";

import React, { useState, useEffect } from 'react';
import { Download, Calendar, HardDrive, FileText, AlertCircle, RefreshCw, HistoryIcon } from 'lucide-react';

interface BackupFile {
  backup_file: string;
  timestamp: string;
  size: number;
  hash: string;
  created_at: string;
  downloadUrl: string;
  githubUrl: string;
  codebergUrl: string;
  codebergRawUrl: string;
  readableDate: string;
  readableSize: string;
}

interface FileBackups {
  lastBackup: string;
  lastHash: string;
  backups: BackupFile[];
}

interface BackupData {
  backups: {
    [filename: string]: FileBackups;
  };
  stats: {
    totalFiles: number;
    totalBackups: number;
    lastBackup: string | null;
    lastBackupTimestamp: string;
  };
  error?: string;
  message?: string;
}

export default function BackupsPage() {
  const [backupData, setBackupData] = useState<BackupData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchBackups = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/backups');
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Yedek verileri yüklenemedi');
      }
      
      setBackupData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Bir hata oluştu');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBackups();
  }, []);

  const handleDownload = (downloadUrl: string, fileName: string) => {
    // Create a temporary link and trigger download from raw URL
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = fileName;
    link.target = '_blank';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
        <main className="max-w-6xl mx-auto px-4 py-16">
          <div className="text-center">
            <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
              Filtre Yedekleri
            </h1>
            <p className="text-lg text-muted-foreground mb-8">
              Turk-AdFilter dosyalarının otomatik yedekleri
            </p>
            <div className="flex justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-red-600"></div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
      <main className="max-w-6xl mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
            Filtre Yedekleri
          </h1>
          <p className="text-lg text-muted-foreground mb-8">
            Turk-AdFilter dosyalarının otomatik yedekleri ve sürüm geçmişi
          </p>
          
          <div className="flex justify-center gap-4 mb-8">
            <button
              onClick={fetchBackups}
              disabled={isLoading}
              className="inline-flex items-center gap-2 px-4 py-2 bg-card border border-border rounded-lg hover:bg-accent hover:text-accent-foreground transition-colors disabled:opacity-50"
            >
              <RefreshCw className={`w-4 h-4 ${isLoading ? 'animate-spin' : ''}`} />
              Yenile
            </button>
          </div>

          {/* Stats */}
          {backupData && !error && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-3xl mx-auto mb-12">
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="flex items-center justify-center mb-2">
                  <FileText className="w-8 h-8 text-red-600 dark:text-red-400" />
                </div>
                <div className="text-2xl font-bold text-foreground mb-1">
                  {backupData.stats.totalFiles}
                </div>
                <p className="text-sm text-muted-foreground">Yedeklenen Dosya</p>
              </div>
              
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="flex items-center justify-center mb-2">
                  <HistoryIcon className="w-8 h-8 text-red-600 dark:text-red-400" />
                </div>
                <div className="text-2xl font-bold text-foreground mb-1">
                  {backupData.stats.totalBackups}
                </div>
                <p className="text-sm text-muted-foreground">Toplam Yedek</p>
              </div>
              
              <div className="bg-card p-6 rounded-lg shadow-sm border">
                <div className="flex items-center justify-center mb-2">
                  <Calendar className="w-8 h-8 text-red-600 dark:text-red-400" />
                </div>
                <div className="text-sm font-bold text-foreground mb-1">
                  {backupData.stats.lastBackup || 'Henüz yok'}
                </div>
                <p className="text-sm text-muted-foreground">Son Yedek</p>
              </div>
            </div>
          )}
        </div>

        {/* Error State */}
        {error && (
          <div className="max-w-2xl mx-auto mb-8">
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
              <div className="flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 mt-0.5" />
                <div>
                  <h3 className="font-medium text-red-800 dark:text-red-200 mb-1">
                    Yedek Verileri Yüklenemedi
                  </h3>
                  <p className="text-red-700 dark:text-red-300 text-sm">
                    {error}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Backup Data with Error */}
        {backupData?.error && (
          <div className="max-w-2xl mx-auto mb-8">
            <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6">
              <div className="flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-yellow-600 dark:text-yellow-400 mt-0.5" />
                <div>
                  <h3 className="font-medium text-yellow-800 dark:text-yellow-200 mb-1">
                    Uyarı
                  </h3>
                  <p className="text-yellow-700 dark:text-yellow-300 text-sm">
                    {backupData.error}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Backup Lists */}
        {backupData && Object.keys(backupData.backups).length > 0 ? (
          <div className="space-y-8">
            {Object.entries(backupData.backups).map(([filename, fileData]) => (
              <div key={filename} className="bg-card border border-border rounded-lg overflow-hidden">
                <div className="bg-muted px-6 py-4 border-b border-border">
                  <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold text-foreground flex items-center gap-2">
                      <FileText className="w-5 h-5" />
                      {filename}
                    </h2>
                    <div className="text-sm text-muted-foreground">
                      {fileData.backups.length} yedek
                    </div>
                  </div>
                  {fileData.lastBackup && (
                    <p className="text-sm text-muted-foreground mt-1">
                      Son yedek: {fileData.lastBackup}
                    </p>
                  )}
                </div>
                
                <div className="divide-y divide-border">
                  {fileData.backups.length === 0 ? (
                    <div className="px-6 py-8 text-center text-muted-foreground">
                      Bu dosya için henüz yedek bulunmuyor
                    </div>
                  ) : (
                    fileData.backups.map((backup, index) => (
                      <div key={backup.backup_file} className="px-6 py-4 hover:bg-muted/50 transition-colors">
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <span className="text-sm font-medium text-foreground">
                                {backup.backup_file}
                              </span>
                              {index === 0 && (
                                <span className="inline-flex items-center px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200 text-xs rounded-md">
                                  En Son
                                </span>
                              )}
                            </div>
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <div className="flex items-center gap-1">
                                <Calendar className="w-4 h-4" />
                                {backup.readableDate}
                              </div>
                              <div className="flex items-center gap-1">
                                <HardDrive className="w-4 h-4" />
                                {backup.readableSize}
                              </div>
                            </div>
                          </div>
                          
                          <div className="flex gap-2">
                            <button
                              onClick={() => handleDownload(backup.downloadUrl, backup.backup_file)}
                              className="inline-flex items-center gap-2 px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm"
                            >
                              <Download className="w-4 h-4" />
                              GitHub
                            </button>
                            <button
                              onClick={() => handleDownload(backup.codebergRawUrl, backup.backup_file)}
                              className="inline-flex items-center gap-2 px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm"
                            >
                              <Download className="w-4 h-4" />
                              Codeberg
                            </button>
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          !error && !backupData?.error && (
            <div className="text-center py-12">
              <div className="max-w-md mx-auto">
                <div className="bg-muted rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                  <HistoryIcon className="w-8 h-8 text-muted-foreground" />
                </div>
                <h3 className="text-lg font-medium text-foreground mb-2">
                  Henüz yedek bulunmuyor
                </h3>
                <p className="text-muted-foreground">
                  Filtre dosyalarında değişiklik yapıldığında otomatik olarak yedek alınacaktır.
                </p>
              </div>
            </div>
          )
        )}

        {/* Info Box */}
        <div className="mt-12 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
          <h3 className="font-medium text-blue-800 dark:text-blue-200 mb-2">
            Yedekleme Sistemi Hakkında
          </h3>
          <ul className="text-blue-700 dark:text-blue-300 text-sm space-y-1">
            <li>• Filtre dosyaları her değişiklikte otomatik olarak yedeklenir</li>
            <li>• Her dosya için en fazla 10 yedek tutulur</li>
            <li>• 30 günden eski yedekler otomatik olarak temizlenir</li>
            <li>• Yedekler SHA256 hash kontrolü ile doğrulanır</li>
            <li>• İndirilen dosyalar doğrudan kullanılabilir</li>
          </ul>
        </div>
      </main>
    </div>
  );
} 