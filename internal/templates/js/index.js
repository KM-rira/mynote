document.querySelectorAll("button").forEach((button) => {
    button.addEventListener("click", () => {
        // ボタンの親要素（行）を取得
        const row = button.parentElement.parentElement;

        // 行内のIDセルを取得
        const id = row.querySelector("td:nth-child(2)").textContent.trim();

        // /select にリダイレクト
        window.location.href = `/select?id=${id}`;
    });
});
