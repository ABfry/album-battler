import Image from "next/image";

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[#f4f4f4]">
      <div className="w-[260px] rounded-lg border border-gray-400 bg-white px-6 py-8 shadow-sm">
        <h1 className="mb-8 text-center text-lg font-semibold">
          アルバムバトラー
        </h1>

        <div className="flex flex-col gap-4">
          <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
            部屋をつくる
          </button>
          <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
            部屋をさがす
          </button>
          <button className="w-full rounded-md border border-gray-400 bg-white py-2 shadow-sm transition hover:-translate-y-0.5 hover:shadow">
            観戦？
          </button>
        </div>
      </div>
    </main>
  );
}
