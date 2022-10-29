export function followLinkPost(url: string) {
    const form = document.createElement("form");
    form.action = url;
    form.method = "POST";
    form.submit();
}
