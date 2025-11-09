import Image from "next/image";
import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[#f4f4f4]">
      <div className="h-[640px] w-[300px] rounded-lg border border-gray-400 bg-white px-6 py-8 shadow-sm">
        <h1 className="mb-8 text-center text-lg font-semibold">
          アルバムバトラー
        </h1>

        <div className="flex flex-col gap-4">
          <Link href="/create">
            <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
              部屋をつくる
            </button>
          </Link>

          <Link href="/search">
            <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
              部屋をさがす
            </button>
          </Link>

          <Link href="/watch">
            <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
              観戦？
            </button>
          </Link>
        </div>
      </div>
    </main>
  );
}
