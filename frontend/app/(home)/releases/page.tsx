"use client";

import React from "react";
import { useState, useEffect } from "react";
import {
  Calendar,
  Tag,
  Download,
  ExternalLink,
  Github,
  ChevronDown,
  Rss,
} from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface GitHubRelease {
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

interface PaginatedResponse {
  releases: GitHubRelease[];
  total: number;
  page: number;
  perPage: number;
  hasMore: boolean;
}

export default function ReleasesPage() {
  const [releases, setReleases] = useState<GitHubRelease[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);

  useEffect(() => {
    const fetchReleases = async () => {
      try {
        const response = await fetch("/api/releases?page=1&per_page=5");
        if (!response.ok) {
          throw new Error("Failed to fetch releases");
        }
        const data: PaginatedResponse = await response.json();
        setReleases(data.releases);
        setHasMore(data.hasMore);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setIsLoading(false);
      }
    };

    fetchReleases();
  }, []);

  const loadMoreReleases = async () => {
    if (isLoadingMore) return;

    setIsLoadingMore(true);
    try {
      const nextPage = currentPage + 1;
      const response = await fetch(`/api/releases?page=${nextPage}&per_page=5`);
      if (!response.ok) {
        throw new Error("Failed to fetch more releases");
      }
      const data: PaginatedResponse = await response.json();

      setReleases((prev) => [...prev, ...data.releases]);
      setCurrentPage(nextPage);
      setHasMore(data.hasMore);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load more releases"
      );
    } finally {
      setIsLoadingMore(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("tr-TR", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const truncateBody = (body: string, maxLength: number = 300) => {
    if (body.length <= maxLength) return body;
    return body.substring(0, maxLength) + "...";
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
        <main className="max-w-6xl mx-auto px-4 py-16">
          <div className="text-center">
            <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
              Sürüm Notları
            </h1>
            <p className="text-lg text-muted-foreground mb-8">
              Turk-AdFilter'ın en son sürümlerini ve değişikliklerini takip edin
            </p>
            <div className="flex justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-red-600"></div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
        <main className="max-w-6xl mx-auto px-4 py-16">
          <div className="text-center">
            <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
              Sürüm Notları
            </h1>
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 max-w-md mx-auto">
              <p className="text-red-800 dark:text-red-200">
                Sürüm notları yüklenirken bir hata oluştu: {error}
              </p>
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
            Sürüm Notları
          </h1>
          <p className="text-lg text-muted-foreground mb-8">
            Turk-AdFilter'ın en son sürümlerini ve değişikliklerini takip edin
          </p>
          <div className="flex justify-center gap-4">
            <a
              href="https://github.com/omerdduran/turk-adfilter/releases"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-4 py-2 bg-card border border-border rounded-lg hover:bg-accent hover:text-accent-foreground transition-colors"
            >
              <Github className="w-4 h-4" />
              GitHub'da Görüntüle
            </a>
            <a
              href="/api/releases/rss"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-lg transition-colors"
            >
              <Rss className="w-4 h-4" />
              RSS Feed
            </a>
          </div>
        </div>

        <div className="space-y-6">
          {releases.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">
                Henüz yayınlanmış sürüm bulunmuyor.
              </p>
            </div>
          ) : (
            releases.map((release) => (
              <div
                key={release.id}
                className="bg-card border border-border rounded-lg p-6 hover:shadow-lg transition-shadow"
              >
                <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-3">
                      <h2 className="text-2xl font-semibold text-foreground">
                        {release.name || release.tag_name}
                      </h2>
                      <span className="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200 text-sm rounded-md">
                        <Tag className="w-3 h-3" />
                        {release.tag_name}
                      </span>
                      {release.prerelease && (
                        <span className="inline-flex items-center px-2 py-1 bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200 text-sm rounded-md">
                          Ön Sürüm
                        </span>
                      )}
                    </div>

                    <div className="flex items-center gap-4 text-sm text-muted-foreground mb-4">
                      <div className="flex items-center gap-1">
                        <Calendar className="w-4 h-4" />
                        {formatDate(release.published_at)}
                      </div>
                      {release.assets.length > 0 && (
                        <div className="flex items-center gap-1">
                          <Download className="w-4 h-4" />
                          {release.assets.length} dosya
                        </div>
                      )}
                    </div>

                    {release.body && (
                      <div className="prose dark:prose-invert max-w-none mb-4">
                        <div className="text-muted-foreground">
                          <ReactMarkdown
                            remarkPlugins={[remarkGfm]}
                            components={{
                              // Customize heading styles
                              h1: ({ children }) => (
                                <h1 className="text-xl font-bold text-foreground mb-3">
                                  {children}
                                </h1>
                              ),
                              h2: ({ children }) => (
                                <h2 className="text-lg font-semibold text-foreground mb-2">
                                  {children}
                                </h2>
                              ),
                              h3: ({ children }) => (
                                <h3 className="text-base font-medium text-foreground mb-2">
                                  {children}
                                </h3>
                              ),
                              h4: ({ children }) => (
                                <h4 className="text-sm font-medium text-foreground mb-1">
                                  {children}
                                </h4>
                              ),
                              h5: ({ children }) => (
                                <h5 className="text-sm font-medium text-foreground mb-1">
                                  {children}
                                </h5>
                              ),
                              h6: ({ children }) => (
                                <h6 className="text-sm font-medium text-foreground mb-1">
                                  {children}
                                </h6>
                              ),
                              // Customize paragraph styles
                              p: ({ children }) => (
                                <p className="mb-2 leading-relaxed">
                                  {children}
                                </p>
                              ),
                              // Customize list styles
                              ul: ({ children }) => (
                                <ul className="list-disc list-inside mb-2 space-y-1">
                                  {children}
                                </ul>
                              ),
                              ol: ({ children }) => (
                                <ol className="list-decimal list-inside mb-2 space-y-1">
                                  {children}
                                </ol>
                              ),
                              li: ({ children }) => (
                                <li className="text-sm">{children}</li>
                              ),
                              // Customize code styles
                              code: ({ children, className }) => {
                                const isInline = !className;
                                return isInline ? (
                                  <code className="bg-muted px-1 py-0.5 rounded text-sm font-mono">
                                    {children}
                                  </code>
                                ) : (
                                  <code className="block bg-muted p-3 rounded text-sm font-mono overflow-x-auto">
                                    {children}
                                  </code>
                                );
                              },
                              pre: ({ children }) => (
                                <pre className="bg-muted p-3 rounded text-sm font-mono overflow-x-auto mb-2">
                                  {children}
                                </pre>
                              ),
                              // Customize link styles
                              a: ({ href, children }) => (
                                <a
                                  href={href}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="text-red-600 dark:text-red-400 hover:underline"
                                >
                                  {children}
                                </a>
                              ),
                              // Customize blockquote styles
                              blockquote: ({ children }) => (
                                <blockquote className="border-l-4 border-red-500 pl-4 italic text-muted-foreground mb-2">
                                  {children}
                                </blockquote>
                              ),
                              // Customize table styles
                              table: ({ children }) => (
                                <div className="overflow-x-auto mb-2">
                                  <table className="min-w-full border border-border rounded">
                                    {children}
                                  </table>
                                </div>
                              ),
                              th: ({ children }) => (
                                <th className="border border-border px-3 py-2 bg-muted font-medium text-left">
                                  {children}
                                </th>
                              ),
                              td: ({ children }) => (
                                <td className="border border-border px-3 py-2">
                                  {children}
                                </td>
                              ),
                            }}
                          >
                            {truncateBody(release.body)}
                          </ReactMarkdown>
                        </div>
                      </div>
                    )}

                    {release.assets.length > 0 && (
                      <div className="mb-4">
                        <h3 className="text-sm font-medium text-foreground mb-2">
                          İndirilebilir Dosyalar:
                        </h3>
                        <div className="flex flex-wrap gap-2">
                          {release.assets.map((asset, index) => (
                            <a
                              key={index}
                              href={asset.browser_download_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="inline-flex items-center gap-2 px-3 py-1 bg-muted hover:bg-accent text-sm rounded-md transition-colors"
                            >
                              <Download className="w-3 h-3" />
                              {asset.name}
                            </a>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="flex flex-col gap-2 lg:flex-shrink-0">
                    <a
                      href={release.html_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center justify-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
                    >
                      <ExternalLink className="w-4 h-4" />
                      Detayları Görüntüle
                    </a>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {hasMore && (
          <div className="mt-8 text-center">
            <button
              onClick={loadMoreReleases}
              disabled={isLoadingMore}
              className="inline-flex items-center gap-2 px-6 py-3 bg-muted hover:bg-accent text-foreground rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoadingMore ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
                  Yükleniyor...
                </>
              ) : (
                <>
                  <ChevronDown className="w-4 h-4" />
                  Daha Fazla Göster
                </>
              )}
            </button>
          </div>
        )}
      </main>
    </div>
  );
}
