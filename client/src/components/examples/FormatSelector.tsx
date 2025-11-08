import { FormatSelector } from "../FormatSelector";
import { useState } from "react";

export default function FormatSelectorExample() {
  const [format, setFormat] = useState("mp4");

  return <FormatSelector selectedFormat={format} onFormatChange={setFormat} />;
}
