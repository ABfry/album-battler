import { Button } from "@/src/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/src/components/ui/card";
import { Input } from "@/src/components/ui/input";
import { Label } from "@/src/components/ui/label";

export default function UiDemoPage() {
  return (
    <div className="min-h-screen p-8">
      <div className="mx-auto max-w-4xl space-y-8">
        {/* タイトル */}
        <div className="text-center">
          <h1 className="text-4xl font-black">アルバムバトラー</h1>
          <p className="mt-2 text-lg">UIコンポーネントデモ</p>
        </div>

        {/* ボタンセクション */}
        <Card>
          <CardHeader>
            <CardTitle>ボタンコンポーネント</CardTitle>
            <CardDescription>
              primaryとsecondaryのボタンスタイル
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-4">
              <Button variant="primary" size="lg" className="w-full">
                primary
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="secondary" size="lg" className="w-full">
                secondary
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="default" size="lg" className="w-full">
                default
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="destructive" size="lg" className="w-full">
                destructive
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="outline" size="lg" className="w-full">
                outline
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="ghost" size="lg" className="w-full">
                ghost
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="link" size="lg" className="w-full">
                link
              </Button>
            </div>
            <div className="flex gap-4">
              <Button variant="primary" size="lg" className="w-full" disabled>
                disabled
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* フォームセクション */}
        <Card>
          <CardHeader>
            <CardTitle>フォームコンポーネント</CardTitle>
            <CardDescription>InputとLabelの組み合わせ</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="number-input">数値入力</Label>
              <Input id="number-input" type="number" placeholder="番号を入力" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="text-input">テキスト入力</Label>
              <Input id="text-input" type="text" placeholder="テキストを入力" />
            </div>
          </CardContent>
        </Card>

        {/* リストカード */}
        <Card>
          <CardHeader>
            <CardTitle>リスト</CardTitle>
            <CardDescription>リスト</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-bold">岩崎</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-bold">井上</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-bold">シバタ</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-bold">なかむら</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-bold">木村</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
