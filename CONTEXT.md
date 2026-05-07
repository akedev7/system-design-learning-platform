# System Design Learning Platform

A Next.js + Go web application for interactive system design learning, inspired by Brilliant's hands-on exercise model.

## Language

### User Roles
**Learner**: An authenticated platform user who enrolls in courses and completes lessons.
_Avoid_: Student, User (when referring to role)

**Admin**: An authenticated platform user with privileges to create, edit, and delete learning content.
_Avoid_: Instructor, Staff

### Learning Content Hierarchy
**Course**: A top-level learning unit focused on a system design topic, containing ordered Modules.
_Avoid_: Class, Track

**Module**: A group of related Lessons within a single Course.
_Avoid_: Section, Chapter

**Lesson**: The smallest unit of learning content, containing ordered ContentBlocks.
_Avoid_: Unit, Page

**ContentBlock**: A typed content element within a Lesson, rendered to Learners.
_Avoid_: Block, Content

### ContentBlock Types
**TextBlock**: A ContentBlock displaying rich text explanations.
_Avoid_: TextContent, ExplanationBlock

**QuizBlock**: A ContentBlock containing interactive questions to test Learner understanding.
_Avoid_: Quiz, TestBlock

**ReactFlowDiagramBlock**: A ContentBlock where Learners build system design diagrams via drag-and-drop.
_Avoid_: DiagramBlock, SystemDesignBlock

**ImageBlock**: A ContentBlock displaying a visual asset from cloud storage.
_Avoid_: Image, PictureBlock

**CodeSnippetBlock**: A ContentBlock displaying formatted code examples.
_Avoid_: CodeBlock, SnippetBlock

### Progress Tracking
**UserCourseProgress**: A record of a Learner's enrollment and completion status for a Course.
_Avoid_: CourseProgress, Enrollment

**UserLessonProgress**: A record of a Learner's completion status and quiz scores for a Lesson.
_Avoid_: LessonProgress, Progress

## Relationships
- A **Learner** enrolls in zero or more **Courses**
- A **Course** contains one or more **Modules**
- A **Module** contains one or more **Lessons**
- A **Lesson** contains zero or more **ContentBlocks** (ordered)
- A **QuizBlock** belongs to exactly one **Lesson**
- A **ReactFlowDiagramBlock** belongs to exactly one **Lesson**
- A **UserCourseProgress** links one **Learner** to one **Course**
- A **UserLessonProgress** links one **Learner** to one **Lesson**

## Example dialogue
> **Dev:** "When a Learner completes all Lessons in a Module, does the Module get marked as complete?"
> **Domain Expert:** "Yes — Module completion is derived from all child Lessons being marked complete by the Learner."
> **Dev:** "Can an Admin be a Learner too?"
> **Domain Expert:** "Yes — Admin role only grants content management privileges, they retain full Learner access."

## Flagged ambiguities
None.
