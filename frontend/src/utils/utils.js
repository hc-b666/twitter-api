export function formatDate(dateStr) {
  if (!dateStr) return null;
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-GB", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function formatDateAndHour(dateStr) {
  if (!dateStr) return null;
  const date = new Date(dateStr);
  return (
    date.toLocaleDateString("en-GB", {
      year: "numeric",
      month: "short",
      day: "numeric",
    }) +
    " " +
    date.toLocaleTimeString("en-GB", {
      hour: "2-digit",
      minute: "2-digit",
    })
  );
}

export const isImage = (post) => {
  return post && post.file_url && post.file_url.match(/\.(jpeg|jpg|gif|png|webp)$/);
};
