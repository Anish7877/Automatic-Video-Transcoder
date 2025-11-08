import { VideoUpload } from "../VideoUpload";
import { useState } from "react";

export default function VideoUploadExample() {
  const [file, setFile] = useState<File | null>(null);

  return <VideoUpload selectedFile={file} onFileSelect={setFile} />;
}
