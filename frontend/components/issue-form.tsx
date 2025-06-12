"use client";

import { useState } from "react";
import Captcha from "./captcha";

export default function IssueForm() {
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    priority: "low",
  });

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [isCaptchaVerified, setIsCaptchaVerified] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    
    if (!isCaptchaVerified) {
      setError("Lütfen CAPTCHA'yı doğrulayın");
      return;
    }

    setIsSubmitting(true);
    setError(null);
    setSuccess(false);

    try {
      const response = await fetch("/api/issues", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          title: formData.title,
          body: `## Açıklama\n${formData.description}\n\n## Önem Seviyesi\n${formData.priority}`,
          labels: ["user-feedback"],
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Bir hata oluştu. Lütfen tekrar deneyin.");
      }

      setSuccess(true);
      setFormData({
        title: "",
        description: "",
        priority: "low",
      });
      setIsCaptchaVerified(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Bir hata oluştu. Lütfen tekrar deneyin.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto p-4">
      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}
      {success && (
        <div className="mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
          Geri bildiriminiz başarıyla gönderildi. Teşekkür ederiz!
        </div>
      )}

      <div className="mb-4">
        <label className="block text-sm font-medium mb-2 text-left">Sorununuz ne?</label>
        <p className="text-left opacity-50 text-xs mb-2">Açıklayıcı bir başlık belirleyin</p>
        <input
          type="text"
          name="title"
          value={formData.title}
          onChange={handleChange}
          required
          className="w-full px-3 py-2 border rounded-md"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium mb-1 text-left">Detaylı Açıklama</label>
        <p className="text-left opacity-50 text-xs mb-2">Sorununuzu ve bu sorunun neden / nasıl oluştuğunu detaylı bir şekilde açıklayın.</p>
        <textarea
          name="description"
          value={formData.description}
          onChange={handleChange}
          required
          rows={4}
          className="w-full px-3 py-2 border rounded-md"
        />
      </div>

      <div className="mb-4">
        <label className="block text-sm font-medium mb-2 text-left">Önem Seviyesi</label>
        <p className="text-left opacity-50 text-xs mb-2">Bu sorun sizin için ne kadar kritik?</p>
        <select
          name="priority"
          value={formData.priority}
          onChange={handleChange}
          className="w-full px-3 py-2 border rounded-md"
        >
          <option value="low">Düşük</option>
          <option value="medium">Orta</option>
          <option value="high">Yüksek</option>
        </select>
      </div>

      <Captcha onVerify={setIsCaptchaVerified} />

      <button
        type="submit"
        disabled={isSubmitting || !isCaptchaVerified}
        className="w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 disabled:opacity-50"
      >
        {isSubmitting ? "Gönderiliyor..." : "Gönder"}
      </button>
    </form>
  );
}
