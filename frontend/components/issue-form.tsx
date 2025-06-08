"use client";

import { useState } from "react";

export default function IssueForm() {
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    priority: "low",
  });

  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  function handleSubmit() {}

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto p-4">
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

      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 disabled:opacity-50"
      >
        {isSubmitting ? "Gönderiliyor..." : "Gönder"}
      </button>
    </form>
  );
}
