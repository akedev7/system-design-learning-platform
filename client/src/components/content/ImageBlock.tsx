import Image from "next/image";
import { ImageConfig } from "@/lib/content-types";

export function ImageBlock({ config }: { config: ImageConfig }) {
  return (
    <figure data-testid="image-block" className="my-6">
      <div className="relative w-full h-64 md:h-96 rounded-lg overflow-hidden">
        <Image
          src={config.src}
          alt={config.alt}
          fill
          className="object-contain"
        />
      </div>
      {config.caption && (
        <figcaption className="mt-2 text-center text-zinc-600 dark:text-zinc-400 text-sm">
          {config.caption}
        </figcaption>
      )}
    </figure>
  );
}