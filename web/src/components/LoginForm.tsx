export const LoginForm = () => {
    return (
        <div className="flex min-h-full flex-col justify-center px-6 lg:px-8 mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
            <form className="flex flex-col gap-4">
                <div className="text-left">
                    <label htmlFor="email" className="block text-sm/6 font-medium text-gray-900 text-left">Email address</label>
                    <div className="mt-2">
                        <input className="block w-full rounded-md border border-gray-300 bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" type="text" required />
                    </div>
                </div>
                <div className="text-left">
                    <label htmlFor="password" className="block text-sm/6 font-medium text-gray-900 text-left">Password</label>
                    <div className="mt-2">
                        <input className="block w-full rounded-md border border-gray-300 bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" type="password" required />
                    </div>
                </div>
                <button className="p-2 bg-blue-500 text-white rounded" type="submit">Sign In</button>
            </form>
        </div>
    )
}
