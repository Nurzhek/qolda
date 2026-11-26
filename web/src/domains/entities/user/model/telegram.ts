// The profile handle may come with @ or as a link: normalize to https://t.me/<nick>
export function telegramLink(value: string) {
  const nick = value.trim().replace(/^https?:\/\/t\.me\//i, '').replace(/^@/, '');
  return nick ? `https://t.me/${encodeURIComponent(nick)}` : '';
}
