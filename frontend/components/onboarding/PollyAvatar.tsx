import Image from 'next/image';

export default function PollyAvatar() {
  return (
    <Image
      src="/random-headshot_cropped.png"
      alt="Lolo AI avatar"
      width={40}
      height={40}
      className="rounded-full border border-border bg-card shrink-0 -translate-y-2 -translate-x-2"
    />
  );
}
