import { LoginForm } from "../components/LoginForm"

export const Login = () => {
    return (
        <div className="flex min-h-screen items-center justify-center bg-gray-50">
            <div className="w-full max-w-md space-y-8 px-4 sm:px-6">
                <div className="flex flex-col items-center">
                    <span className="text-5xl">🐗</span>
                    <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">Please sign in</h2>
                    <LoginForm />
                </div>
            </div>
        </div>
    )
}
