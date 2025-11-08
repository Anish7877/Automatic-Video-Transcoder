# Design Guidelines: Automatic Video Transcoder

## Design Philosophy
Modern, elegant, and minimal aesthetic inspired by Figma, Notion, and Linear.app. Clean, flat design with smooth interactions and professional polish.

## Layout & Structure

**Single-Page Centered Layout**
- Vertically and horizontally centered main card/panel
- Clean header with app name "Automatic Video Transcoder"
- Generous spacing with consistent grid alignment
- Rounded corners throughout
- Subtle shadows for depth

## Color Palette

**Background:** Soft neutral (#f5f6f8 or #fafafa)  
**Accent Color:** Choose one strong color - #007AFF, #00ADB5, or #6366F1  
**Text:** Dark gray with subtle contrast for labels  
**Success States:** Appropriate green tones  
**Progress/Loading:** Use accent color

## Typography

**Font Family:** Inter, Poppins, Nunito Sans, or DM Sans (sans-serif)  
**Hierarchy:**
- Header: Large, bold app title
- Section Labels: Medium weight, clear hierarchy
- Body Text: Regular weight for file details and instructions
- Buttons: Medium/semi-bold for clarity

## Component Design

### Upload Section
- Prominent drag-and-drop area with dashed border
- Clear visual feedback on drag-over (accent color highlight)
- Fallback "Browse File" button inside upload area
- File metadata display after upload: name, type, size, mock duration/resolution
- Upload icon and instructional text when empty

### Format Selection
- Dropdown or segmented button group
- Formats: .mp4, .avi, .mov, .webm, .mkv
- Single selection only
- Clear selected state with accent color

### Conversion Section
- Prominent "Convert" button using accent color
- Animated progress bar during conversion (smooth, linear animation)
- Alternative: Circular loader with "Converting..." text
- Progress indicator should feel fluid and professional

### Download Section
- "Download Converted File" button revealed after success
- Success feedback with checkmark or completion icon
- Clear visual state transition from converting to complete

### Toast Notifications
- Subtle, non-intrusive toasts for:
  - "Upload successful"
  - "Conversion started"  
  - "Conversion complete"
- Auto-dismiss after 3-4 seconds
- Position: Top-right or top-center

## Spacing System
Use Tailwind spacing: p-4, p-6, p-8 for consistent padding  
Gap between sections: 6-8 units  
Card padding: 8-10 units for breathing room

## Animations & Interactions

**Smooth Transitions:**
- Button hover states with scale or brightness change
- Progress bar smooth fill animation
- State changes fade in/out gracefully
- File upload drag-over visual feedback

**Motion Principles:**
- Subtle, purposeful animations (using Framer Motion or CSS)
- Duration: 200-300ms for most interactions
- Easing: Use smooth curves, avoid linear

## Dark/Light Mode Toggle
- Subtle toggle in header or corner
- Smooth theme transition
- Maintain color contrast ratios
- Persist user preference

## Responsive Behavior

**Desktop:** Full centered card layout with generous spacing  
**Tablet:** Slightly condensed, maintains card structure  
**Mobile:** Vertical stack, full-width sections with appropriate padding

## Accessibility
- ARIA labels for all interactive elements
- Clear focus states with accent color outline
- Keyboard navigation support
- Sufficient color contrast for text
- Screen reader friendly file metadata

## Footer
Small branding text or note ("Powered by [Your Name]")  
Subtle, minimal presence

## Images
No hero images required. Focus on clean UI with icon-based visuals:
- Upload icon in drag-drop area
- Format icons for file types
- Success/completion icons for states
- Loading/spinner animations

## Key Principles
- Simplicity over complexity
- Smooth, natural flow: upload → convert → download
- Professional polish in every interaction
- Generous whitespace for breathing room
- Consistent visual language throughout