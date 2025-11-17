type CreateUserRequest = {
  user_name: string;
};

type CreateUserResponse = {
  user_id: string;
};

export async function createUser(userName: string): Promise<string> {
  if (!process.env.NEXT_PUBLIC_API_URL) {
    throw new Error("API URLが設定されていません");
  }
  const url = `${process.env.NEXT_PUBLIC_API_URL}/user`;
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_name: userName } satisfies CreateUserRequest),
  });

  if (!response.ok) {
    throw new Error("ユーザー作成に失敗しました");
  }

  const data: CreateUserResponse = await response.json();
  return data.user_id;
}
