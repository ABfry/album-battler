"use client";

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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/src/components/ui/alert-dialog";

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

        {/* アラートダイアログセクション */}
        <Card>
          <CardHeader>
            <CardTitle>アラートダイアログ</CardTitle>
            <CardDescription>確認ダイアログ（はい/いいえ）</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* 基本的な確認ダイアログ */}
            <div>
              <p className="mb-2 text-sm font-medium">基本的な確認</p>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="primary">確認ダイアログを開く</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>確認</AlertDialogTitle>
                    <AlertDialogDescription>
                      この操作を実行してもよろしいですか？
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>いいえ</AlertDialogCancel>
                    <AlertDialogAction>はい</AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>

            {/* 削除確認ダイアログ */}
            <div>
              <p className="mb-2 text-sm font-medium">削除確認</p>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">削除する</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>本当に削除しますか？</AlertDialogTitle>
                    <AlertDialogDescription>
                      この操作は取り消せません。本当に削除してもよろしいですか？
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>キャンセル</AlertDialogCancel>
                    <AlertDialogAction
                      className="bg-destructive hover:bg-destructive/90 text-white"
                      onClick={() => alert("削除されました")}
                    >
                      削除
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
