import { Permission } from "@/config/enums";

export function formatSize(size: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let index = 0;
  while (size >= 1024 && index < units.length - 1) {
    size /= 1024;
    index += 1;
  }
  return `${size.toFixed(2)} ${units[index]}`;
}

export function calculatePublic(parentPublic: boolean, permission: Permission) {
  if (permission === Permission.inherit) {
    return parentPublic;
  }
  return permission === Permission.public;
}

export function downloadFile(url: string, name: string) {
  const a = document.createElement('a');
  a.href = url;
  a.download = name;
  a.click(); 
}
