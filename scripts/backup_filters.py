#!/usr/bin/env python3
"""
Otomatik Filtre Yedekleme Sistemi
Turk-AdFilter projesi için hosts.txt ve turk-adfilter.txt dosyalarının yedeklenmesi
"""

import os
import shutil
import datetime
import hashlib
import json
import argparse
from pathlib import Path

class FilterBackup:
    def __init__(self, backup_dir="yedekler"):
        self.backup_dir = Path(backup_dir)
        self.filter_files = ["hosts.txt", "turk-adfilter.txt"]
        self.metadata_file = self.backup_dir / "backup_metadata.json"
        
        # Yedek klasörünü oluştur
        self.backup_dir.mkdir(exist_ok=True)
        
    def get_file_hash(self, filepath):
        """Dosyanın SHA256 hash değerini hesapla"""
        hasher = hashlib.sha256()
        try:
            with open(filepath, 'rb') as f:
                for chunk in iter(lambda: f.read(4096), b""):
                    hasher.update(chunk)
            return hasher.hexdigest()
        except FileNotFoundError:
            return None
    
    def load_metadata(self):
        """Yedek metadata'sını yükle"""
        if self.metadata_file.exists():
            try:
                with open(self.metadata_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            except (json.JSONDecodeError, FileNotFoundError):
                return {}
        return {}
    
    def save_metadata(self, metadata):
        """Yedek metadata'sını kaydet"""
        with open(self.metadata_file, 'w', encoding='utf-8') as f:
            json.dump(metadata, f, indent=2, ensure_ascii=False)
    
    def delete_backup_file(self, backup):
        """Yedek dosyasını diskten sil"""
        backup_path = self.backup_dir / backup['backup_file']
        if backup_path.exists():
            backup_path.unlink()
            return True
        return False

    def dedupe_backups(self, quiet=False):
        """Aynı hash'e sahip yedeklerden yalnızca en yenisini tut."""
        metadata = self.load_metadata()
        deleted = 0

        for file_meta in metadata.values():
            backups = file_meta.get('backups', [])
            seen = {}
            kept = []
            removed = []
            for backup in backups:
                previous = seen.get(backup.get('hash'))
                if previous is not None:
                    kept.remove(previous)
                    removed.append(previous)
                seen[backup.get('hash')] = backup
                kept.append(backup)
            if kept:
                referenced = {backup['backup_file'] for backup in kept}
                for old_backup in removed:
                    if old_backup['backup_file'] not in referenced:
                        self.delete_backup_file(old_backup)
                    deleted += 1
                file_meta['backups'] = kept
                file_meta['last_backup'] = kept[-1]['timestamp']

        if deleted:
            self.save_metadata(metadata)
            print(f"🧹 {deleted} tekrar eden yedek silindi")
        elif not quiet:
            print("✅ Tekrar eden yedek yok")
        return deleted

    def needs_backup(self, filename):
        """Dosyanın yedeklenmesi gerekip gerekmediğini kontrol et"""
        if not os.path.exists(filename):
            return False
            
        current_hash = self.get_file_hash(filename)
        metadata = self.load_metadata()
        
        # İlk yedekleme veya hash değişmişse yedekle
        if filename not in metadata or metadata[filename].get('last_hash') != current_hash:
            return True
        return False
    
    def create_backup(self, filename, force=False):
        """Belirtilen dosyanın yedeğini oluştur"""
        if not os.path.exists(filename):
            print(f"❌ Dosya bulunamadı: {filename}")
            return False
        
        if not force and not self.needs_backup(filename):
            print(f"✅ {filename} için yedek gerekli değil (değişiklik yok)")
            return False
        
        # Zaman damgası oluştur
        timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
        
        # Yedek dosya adını oluştur
        file_ext = Path(filename).suffix
        file_stem = Path(filename).stem
        backup_filename = f"{file_stem}_{timestamp}{file_ext}"
        backup_path = self.backup_dir / backup_filename
        
        try:
            # Dosyayı kopyala
            shutil.copy2(filename, backup_path)
            
            # Metadata güncelle
            metadata = self.load_metadata()
            current_hash = self.get_file_hash(filename)
            
            if filename not in metadata:
                metadata[filename] = {'backups': []}
            
            backup_info = {
                'backup_file': backup_filename,
                'timestamp': timestamp,
                'size': os.path.getsize(filename),
                'hash': current_hash,
                'created_at': datetime.datetime.now().isoformat()
            }
            
            metadata[filename]['backups'].append(backup_info)
            metadata[filename]['last_hash'] = current_hash
            metadata[filename]['last_backup'] = timestamp
            
            # Eski yedekleri temizle (sadece son 10 yedek)
            if len(metadata[filename]['backups']) > 10:
                old_backups = metadata[filename]['backups'][:-10]
                metadata[filename]['backups'] = metadata[filename]['backups'][-10:]
                
                # Eski yedek dosyalarını sil
                for old_backup in old_backups:
                    if self.delete_backup_file(old_backup):
                        print(f"🗑️  Eski yedek silindi: {old_backup['backup_file']}")
            
            self.save_metadata(metadata)

            # Aynı içerikli eski kopyaları temizle (ör. --force tekrarları)
            self.dedupe_backups(quiet=True)
            
            print(f"✅ Yedek oluşturuldu: {filename} -> {backup_filename}")
            return True
            
        except Exception as e:
            print(f"❌ Yedekleme hatası {filename}: {e}")
            return False
    
    def backup_all_filters(self, force=False):
        """Tüm filtre dosyalarını yedekle"""
        success_count = 0
        
        print("🔄 Filtre dosyaları yedekleniyor...")
        
        for filename in self.filter_files:
            if self.create_backup(filename, force):
                success_count += 1
        
        print(f"✅ Tamamlandı: {success_count}/{len(self.filter_files)} dosya yedeklendi")
        return success_count > 0
    
    def list_backups(self, filename=None):
        """Yedekleri listele"""
        metadata = self.load_metadata()
        
        if filename:
            files_to_show = [filename] if filename in metadata else []
        else:
            files_to_show = list(metadata.keys())
        
        if not files_to_show:
            print("📁 Yedek bulunamadı")
            return
        
        for file in files_to_show:
            print(f"\n📄 {file}:")
            if 'backups' in metadata[file]:
                for i, backup in enumerate(reversed(metadata[file]['backups'])):
                    print(f"  {i+1}. {backup['backup_file']} ({backup['timestamp']})")
            else:
                print("  Yedek bulunamadı")
    
    def get_backup_info(self):
        """Yedek bilgilerini al"""
        metadata = self.load_metadata()
        info = {
            'total_files': len(metadata),
            'total_backups': sum(len(meta.get('backups', [])) for meta in metadata.values()),
            'last_backup': None
        }
        
        # En son yedek zamanını bul
        all_backups = []
        for file_meta in metadata.values():
            all_backups.extend(file_meta.get('backups', []))
        
        if all_backups:
            latest = max(all_backups, key=lambda x: x['timestamp'])
            info['last_backup'] = latest['timestamp']
        
        return info

def main():
    parser = argparse.ArgumentParser(description='Filtre dosyalarını yedekle')
    parser.add_argument('--force', action='store_true', help='Değişiklik olmasa da zorla yedekle')
    parser.add_argument('--list', action='store_true', help='Mevcut yedekleri listele')
    parser.add_argument('--file', help='Belirli bir dosyayı yedekle')
    parser.add_argument('--info', action='store_true', help='Yedek bilgilerini göster')
    parser.add_argument('--dedupe', action='store_true', help='Aynı içerikli yedeklerden yalnızca en yenisini tut')
    
    args = parser.parse_args()
    
    backup = FilterBackup()
    
    if args.dedupe:
        backup.dedupe_backups()
        
    elif args.info:
        info = backup.get_backup_info()
        print(f"📊 Yedek İstatistikleri:")
        print(f"   Toplam dosya: {info['total_files']}")
        print(f"   Toplam yedek: {info['total_backups']}")
        print(f"   Son yedek: {info['last_backup'] or 'Yok'}")
        
    elif args.list:
        backup.list_backups(args.file)
        
    elif args.file:
        backup.create_backup(args.file, args.force)
        
    else:
        backup.backup_all_filters(args.force)

if __name__ == "__main__":
    main() 