"use client";

import { useState, useEffect } from "react";

interface CaptchaProps {
  onVerify: (isValid: boolean) => void;
}

export default function Captcha({ onVerify }: CaptchaProps) {
  const [captchaId, setCaptchaId] = useState<string | null>(null);
  const [captchaImage, setCaptchaImage] = useState<string | null>(null);
  const [userAnswer, setUserAnswer] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);
  const [isVerified, setIsVerified] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const generateCaptcha = async () => {
    setIsGenerating(true);
    setError(null);
    
    try {
      const response = await fetch("/api/captcha", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          action: "generate",
          level: "easy",
          media: "image/png",
          input_type: "text",
          size: "350x100",
        }),
      });

      if (!response.ok) {
        throw new Error("CAPTCHA oluşturulamadı");
      }

      const data = await response.json();
      setCaptchaId(data.id);
      
      // Use our proxy for the image
      const imageUrl = `/api/captcha?id=${data.id}`;
      setCaptchaImage(imageUrl);
    } catch (err) {
      setError(err instanceof Error ? err.message : "CAPTCHA yüklenirken hata oluştu");
    } finally {
      setIsGenerating(false);
    }
  };

  const verifyCaptcha = async () => {
    if (!captchaId || !userAnswer.trim()) {
      setError("Lütfen CAPTCHA'yı çözün");
      return;
    }

    setIsVerifying(true);
    setError(null);

    try {
      const response = await fetch("/api/captcha", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          action: "verify",
          id: captchaId,
          answer: userAnswer,
        }),
      });

      if (!response.ok) {
        throw new Error("CAPTCHA doğrulanamadı");
      }

      const data = await response.json();
      
      if (data.result === "True") {
        setIsVerified(true);
        onVerify(true);
        setError(null);
      } else if (data.result === "Expired") {
        setError("CAPTCHA süresi doldu. Lütfen yeni bir CAPTCHA oluşturun.");
        setIsVerified(false);
        onVerify(false);
      } else {
        setError("Yanlış cevap. Lütfen tekrar deneyin.");
        setIsVerified(false);
        onVerify(false);
        setUserAnswer("");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "CAPTCHA doğrulanırken hata oluştu");
      setIsVerified(false);
      onVerify(false);
    } finally {
      setIsVerifying(false);
    }
  };

  useEffect(() => {
    generateCaptcha();
  }, []);

  return (
    <div className="mb-4">
      <label className="block text-sm font-medium mb-1 text-left">
        Güvenlik Doğrulaması
      </label>
      <p className="text-left opacity-50 text-xs mb-1">
        Lütfen aşağıdaki kodu girin
      </p>

      <div className="flex flex-col space-y-3">
        {captchaImage && (
          <div className="flex justify-center">
            <img
              src={captchaImage}
              alt="CAPTCHA"
              className="border rounded"
              style={{ maxWidth: "350px", height: "100px" }}
              onError={(e) => {
                setError("CAPTCHA resmi yüklenemedi");
              }}
              onLoad={() => {
                // Image loaded successfully
              }}
            />
          </div>
        )}

        {isGenerating && !captchaImage && (
          <div className="flex justify-center">
            <div className="border rounded bg-gray-100 flex items-center justify-center" style={{ width: "350px", height: "100px" }}>
              CAPTCHA yükleniyor...
            </div>
          </div>
        )}

        <div className="flex space-x-2">
          <input
            type="text"
            value={userAnswer}
            onChange={(e) => setUserAnswer(e.target.value)}
            placeholder="CAPTCHA kodunu girin"
            disabled={isVerified}
            className="flex-1 px-3 py-2 border rounded-md disabled:opacity-50"
          />
          
          {!isVerified && (
            <button
              type="button"
              onClick={verifyCaptcha}
              disabled={isVerifying || isGenerating || !userAnswer.trim()}
              className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50 min-w-[80px] flex items-center justify-center"
            >
              {isVerifying ? (
                <svg 
                  className="w-4 h-4 animate-spin" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <circle 
                    cx="12" 
                    cy="12" 
                    r="10" 
                    stroke="currentColor" 
                    strokeWidth="4"
                    fill="none"
                    strokeDasharray="31.416"
                    strokeDashoffset="31.416"
                    className="opacity-25"
                  />
                  <path 
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    fill="currentColor"
                    className="opacity-75"
                  />
                </svg>
              ) : (
                "Doğrula"
              )}
            </button>
          )}

          {!isVerified && (
            <button
              type="button"
              onClick={generateCaptcha}
              disabled={isGenerating || isVerifying}
              className="px-3 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 disabled:opacity-50"
              title="Yeni CAPTCHA"
            >
              <svg 
                className={`w-4 h-4 ${isGenerating ? 'animate-spin' : ''}`} 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path 
                  strokeLinecap="round" 
                  strokeLinejoin="round" 
                  strokeWidth={2} 
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" 
                />
              </svg>
            </button>
          )}
        </div>

        {error && (
          <div className="p-2 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
            {error}
          </div>
        )}

        {isVerified && (
          <div className="p-2 bg-green-100 border border-green-400 text-green-700 rounded text-sm">
            ✓ CAPTCHA doğrulandı
          </div>
        )}
      </div>
    </div>
  );
} 