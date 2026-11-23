#!/usr/bin/env python3
"""
Filtre Yedekleri Geri Yükleme Sistemi
Turk-AdFilter projesi için yedeklerden geri yükleme işlemleri
"""

import os
import shutil
import json
import argparse
from pathlib import Path
from datetime import datetime

class FilterRestore:
    def __init__(self, backup_dir="yedekler"):
        self.backup_dir = Path(backup_dir)
        self.metadata_file = self.backup_dir / "backup_metadata.json"

    def load_metadata(self):
        """Yedek metadata'sını yükle"""
        if self.metadata_file.exists():
            try:
                with open(self.metadata_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            except (json.JSONDecodeError, FileNotFoundError):
                return {}
        return {}

    def list_available_backups(self, filename=None):
        """Mevcut yedekleri listele"""
        metadata = self.load_metadata()

        if not metadata:
            print("📁 Yedek bulunamadı")
            return {}

        if filename:
            if filename not in metadata:
                print(f"❌ {filename} için yedek bulunamadı")
                return {}
            files_to_show = {filename: metadata[filename]}
        else:
            files_to_show = metadata

        backup_options = {}
        option_counter = 1

        for file, file_meta in files_to_show.items():
            print(f"\n📄 {file}:")

            if 'backups' not in file_meta or not file_meta['backups']:
                print("  Yedek bulunamadı")
                continue

            # Yedekleri en yeniden en eskiye sırala
            sorted_backups = sorted(file_meta['backups'],
                                  key=lambda x: x['timestamp'],
                                  reverse=True)

            for backup in sorted_backups:
                # Dosya boyutunu MB/KB cinsinden göster
                size_bytes = backup.get('size', 0)
                if size_bytes > 1024*1024:
                    size_str = f"{size_bytes/(1024*1024):.1f} MB"
                elif size_bytes > 1024:
                    size_str = f"{size_bytes/1024:.1f} KB"
                else:
                    size_str = f"{size_bytes} B"

                # Tarihi readable formata çevir
                try:
                    timestamp = backup['timestamp']
                    dt = datetime.strptime(timestamp, "%Y%m%d_%H%M%S")
                    readable_date = dt.strftime("%d.%m.%Y %H:%M:%S")
                except Exception:
                    readable_date = timestamp

                print(f"  {option_counter}. {backup['backup_file']}")
                print(f"     📅 {readable_date} | 📦 {size_str}")

                backup_options[option_counter] = {
                    'original_file': file,
                    'backup_file': backup['backup_file'],
                    'timestamp': timestamp,
                    'readable_date': readable_date,
                    'size': size_str
                }
                option_counter += 1

        return backup_options

    def restore_backup(self, backup_option, confirm=True):
        """Seçilen yedekten geri yükleme yap"""
        original_file = backup_option['original_file']
        backup_file = backup_option['backup_file']
        backup_path = self.backup_dir / backup_file

        if not backup_path.exists():
            print(f"❌ Yedek dosya bulunamadı: {backup_file}")
            return False

        if confirm:
            print("\n🔄 Geri Yükleme Onayı:")
            print(f"   Dosya: {original_file}")
            print(f"   Yedek: {backup_file}")
            print(f"   Tarih: {backup_option['readable_date']}")
            print(f"   Boyut: {backup_option['size']}")

            response = input("\nGeri yüklemeyi onaylıyor musunuz? (e/H): ").lower()
            if response not in ['e', 'evet', 'y', 'yes']:
                print("❌ İşlem iptal edildi")
                return False

        try:
            # Mevcut dosyanın yedeğini al (güvenlik için)
            if os.path.exists(original_file):
                safety_backup = f"{original_file}.restore_safety_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
                shutil.copy2(original_file, safety_backup)
                print(f"🛡️  Güvenlik yedeği oluşturuldu: {safety_backup}")

            # Yedekten geri yükle
            shutil.copy2(backup_path, original_file)

            print(f"✅ Geri yükleme başarılı: {original_file}")
            print(f"   Kaynak: {backup_file}")
            return True

        except Exception as e:
            print(f"❌ Geri yükleme hatası: {e}")
            return False

    def interactive_restore(self):
        """İnteraktif geri yükleme menüsü"""
        print("🔄 Filtre Geri Yükleme Sistemi")
        print("=" * 40)

        backup_options = self.list_available_backups()

        if not backup_options:
            return False

        print(f"\n📋 Geri yüklenebilir {len(backup_options)} yedek bulundu")
        print("\nLütfen geri yüklemek istediğiniz yedek numarasını girin:")

        try:
            choice = int(input("\nSeçim (0=çıkış): "))

            if choice == 0:
                print("❌ İşlem iptal edildi")
                return False

            if choice not in backup_options:
                print(f"❌ Geçersiz seçim: {choice}")
                return False

            return self.restore_backup(backup_options[choice])

        except ValueError:
            print("❌ Geçersiz giriş")
            return False
        except KeyboardInterrupt:
            print("\n❌ İşlem iptal edildi")
            return False

    def restore_latest(self, filename):
        """Belirtilen dosyanın en son yedeğini geri yükle"""
        metadata = self.load_metadata()

        if filename not in metadata or not metadata[filename].get('backups'):
            print(f"❌ {filename} için yedek bulunamadı")
            return False

        # En son yedek
        latest_backup = max(metadata[filename]['backups'],
                          key=lambda x: x['timestamp'])

        backup_option = {
            'original_file': filename,
            'backup_file': latest_backup['backup_file'],
            'timestamp': latest_backup['timestamp'],
            'readable_date': latest_backup['timestamp'],
            'size': f"{latest_backup.get('size', 0)} B"
        }

        return self.restore_backup(backup_option, confirm=False)

    def cleanup_old_backups(self, days=30):
        """Belirtilen günden eski yedekleri temizle"""
        metadata = self.load_metadata()
        current_time = datetime.now()
        cleaned_count = 0

        for filename, file_meta in metadata.items():
            if 'backups' not in file_meta:
                continue

            backups_to_keep = []

            for backup in file_meta['backups']:
                try:
                    backup_time = datetime.strptime(backup['timestamp'], "%Y%m%d_%H%M%S")
                    age_days = (current_time - backup_time).days

                    if age_days <= days:
                        backups_to_keep.append(backup)
                    else:
                        # Eski yedek dosyasını sil
                        backup_path = self.backup_dir / backup['backup_file']
                        if backup_path.exists():
                            backup_path.unlink()
                            cleaned_count += 1
                            print(f"🗑️  Eski yedek silindi: {backup['backup_file']} ({age_days} gün)")
                except Exception:
                    # Hatalı timestamp varsa koru
                    backups_to_keep.append(backup)

            metadata[filename]['backups'] = backups_to_keep

        # Güncellenen metadata'yı kaydet
        if cleaned_count > 0:
            with open(self.metadata_file, 'w', encoding='utf-8') as f:
                json.dump(metadata, f, indent=2, ensure_ascii=False)

        print(f"✅ Temizlik tamamlandı: {cleaned_count} eski yedek silindi")
        return cleaned_count

def main():
    parser = argparse.ArgumentParser(description='Filtre yedeklerini geri yükle')
    parser.add_argument('--list', action='store_true', help='Mevcut yedekleri listele')
    parser.add_argument('--file', help='Belirli bir dosyanın yedeklerini listele')
    parser.add_argument('--restore-latest', help='Belirtilen dosyanın en son yedeğini geri yükle')
    parser.add_argument('--interactive', action='store_true', help='İnteraktif geri yükleme menüsü')
    parser.add_argument('--cleanup', type=int, metavar='DAYS', help='Belirtilen günden eski yedekleri temizle')

    args = parser.parse_args()

    restore = FilterRestore()

    if args.cleanup:
        restore.cleanup_old_backups(args.cleanup)

    elif args.list:
        restore.list_available_backups(args.file)

    elif args.restore_latest:
        restore.restore_latest(args.restore_latest)

    elif args.interactive:
        restore.interactive_restore()

    else:
        # Varsayılan olarak interaktif menüyü göster
        restore.interactive_restore()

if __name__ == "__main__":
    main()
