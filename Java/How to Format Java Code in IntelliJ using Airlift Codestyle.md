In this blog post, I'll guide you through the process of setting up and using a predefined Java code style in IntelliJ, specifically using the Airlift code style configuration for IntelliJ IDEA 2019 or newer.

### Why Use a Predefined Code Style?

Before we dive into the how-to, let's briefly discuss why using a predefined code style can be beneficial:

1.  **Consistency:** A consistent coding style helps maintain the readability and maintainability of your code, which is crucial for team projects or when your projects grow in size.
2.  **Efficiency:** With predefined styles, you can automate formatting, saving you time and effort in making your code look tidy.
3.  **Best Practices:** Often, these styles are designed keeping best practices in mind, helping you adhere to widely accepted coding standards.

> **NOTE**: Make sure you're using IntelliJ IDEA 2019 or newer to follow along with these instructions.

### Step 1: Download the Airlift Code Style File

The Airlift code style is a set of configurations optimized for Java that you can directly import into IntelliJ IDEA. To start, download the `Airlift.xml` file from the Airlift GitHub repository. Here is the link for your convenience:

### [Download Airlift Code Style](https://github.com/airlift/codestyle/blob/master/IntelliJIdea2019/Airlift.xml)

### Step 2: Import the Airlift Scheme into IntelliJ

Once you have the `Airlift.xml` file, you need to import it into your IntelliJ IDEA environment. Follow these steps:

1.  Open IntelliJ IDEA and go to `Preferences` or `sittings`. On Windows and Linux, you can access this by going to `File > Settings`. On macOS, it's under `IntelliJ IDEA > Preferences`.
2.  In the Preferences window, navigate to `Editor > Code Style`.
3.  Click on the gear icon (settings) next to the Scheme dropdown at the top of the Code Style window.
4.  Select `Import Scheme` and then `IntelliJ IDEA code style XML`.
5.  Navigate to where you saved the `Airlift.xml` file, select it, and click OK.
6.  Once imported, make sure that Airlift is selected as the current code style scheme.

   ![instruction](https://github.com/code-flu/website-content/assets/104217888/28a24be0-8818-4ce0-ab3e-50fd0156225a)


### Step 4: Use the Airlift Code Style in Your Projects

With the Airlift scheme now active, IntelliJ IDEA will automatically format your Java code according to the rules specified in the Airlift style guide whenever you write code or reformat existing code. You can manually reformat your code by using the shortcut `Ctrl+Alt+L` on Windows/Linux or `Cmd+Alt+L` on macOS.

### Conclusion

Happy coding, and enjoy the streamlined beauty of a well-formatted codebase in IntelliJ IDEA!
